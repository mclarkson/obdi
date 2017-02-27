// Copyright 2016 Mark Clarkson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"time"
	//"encoding/json"
	//"strings"
)

func (api *Api) execCmd(job JobIn) {

	// TODO :: Put this logic in login/logout and reference count
	//defer api.Logout( )

	// Need to set the PATH to run the script from the script dir
	os.Setenv("PATH", config.ScriptDir)

	scriptfile := ""

	// Write ScriptSource to disk
	if file, err := ioutil.TempFile(os.TempDir(), "smworker_"); err != nil {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_SYSCANCELLED,
			StatusReason:  fmt.Sprintf("TempFile error ('%s')", err.Error()),
			StatusPercent: 0,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
		return
	} else {
		if _, err := file.Write(job.ScriptSource); err != nil {
			if err := api.sendStatus(job, JobOut{
				Status:        STATUS_SYSCANCELLED,
				StatusReason:  fmt.Sprintf("Write error ('%s')", err.Error()),
				StatusPercent: 0,
				Errors:        0,
			}); err != nil {
				logit(fmt.Sprintf("Error: %s", err.Error()))
			}
			return
		}
		file.Close()
		os.Chmod(file.Name(), 0755)
		scriptfile = file.Name()
	}
	defer os.Remove(scriptfile)

	// Set up command, split on spaces but preserve quoted strings
	head := scriptfile
	r := regexp.MustCompile("'.+'|\".+\"|\\S+")
	parts := r.FindAllString(job.Args, -1)
	cmd := &exec.Cmd{}
	cmd = exec.Command(head, parts...)

	// Apply the sent environment variables, split on spaces
	// but preserve quoted strings
	if len(job.EnvVars) > 0 {
		r = regexp.MustCompile(`[^ ]*='[^']+'|[^ ]*="[^"]+"|\S+`)
		cmd.Env = r.FindAllString(job.EnvVars, -1)
		// Remove speech marks around the value of quoted strings. Matches, for
		// example, `var="val val"`, and changes to `var=val val`
		for i, j := range cmd.Env {
			r := regexp.MustCompile(`(?sU)([^ ]+=)["'](.*)["']`)
			if r.Match([]byte(j)) {
				k := r.ReplaceAll([]byte(j), []byte(`$1$2`))
				cmd.Env[i] = string(k)
			}
		}
	} else {
		cmd.Env = []string{}
	}

	// Add the system scripts directory to Env.SYSSCRIPTDIR
	cmd.Env = append(cmd.Env, "SYSSCRIPTDIR="+config.SysScriptDir)

	cmd.Dir = os.TempDir()

	// Set up buffer for stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_SYSCANCELLED,
			StatusReason:  fmt.Sprintf("Pipe error ('%s')", err.Error()),
			StatusPercent: 0,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
		return
	}
	rdr := bufio.NewReader(stdout)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_SYSCANCELLED,
			StatusReason:  fmt.Sprintf("Pipe error ('%s')", err.Error()),
			StatusPercent: 0,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
		return
	}
	rdr_stderr := bufio.NewReader(stderr)

	/*
	   data := JobOut{}
	   jsondata, err := json.Marshal(data)
	*/

	// Get child processes to run in a process group so they
	// can all be killed as a group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	// Run command in the background (fork)
	err = cmd.Start()
	if err != nil {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_SYSCANCELLED,
			StatusReason:  fmt.Sprintf("Fork error ('%s')", err.Error()),
			StatusPercent: 0,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
		return
	}

	if err := api.sendStatus(job, JobOut{
		Status:        STATUS_INPROGRESS,
		StatusReason:  "Script started",
		StatusPercent: 0,
		Errors:        0,
	}); err != nil {
		logit(fmt.Sprintf("Error: %s", err.Error()))
	}

	// Save the pid so it can be killed
	api.SetPid(job.JobID, int64(cmd.Process.Pid))

	// Process the output
	serial := int64(1)
	line := ""
	if job.Type != 2 {
		input := make(chan string, 1)

		getInput := func(input chan string) {
			for {
				result, err := rdr.ReadString('\n')
				input <- result
				if err != nil {
					close(input)
					break
				}
			}
		}

		go getInput(input)
		lines := ""
		t := time.NewTicker(time.Duration(1000 * time.Millisecond))

	OuterLoop:
		for {
			select {
			case i, ok := <-input:
				lines += i
				if ok == false {
					api.sendOutputLine(job, lines, serial)
					serial++
					break OuterLoop
				}
			case <-t.C:
				if len(lines) > 0 {
					api.sendOutputLine(job, lines, serial)
					serial++
					lines = ""
				}
			}
		}

		t.Stop()

	} else {
		// A system job. Send all output in a single output line
		a := ""
		for err == nil {
			line, err = rdr.ReadString('\n')
			a = a + line
		}
		api.sendOutputLine(job, a, 1)
	}

	serial++

	// Read anything in stderr, but don't send it.
	// It will be sent later if the script has non-zero exit status
	error_output := ""
	err = nil
	for err == nil {
		line, err = rdr_stderr.ReadString('\n')
		error_output = error_output + line
	}

	// Process exit status
	err = cmd.Wait()
	if err != nil {
		api.sendOutputLine(job, error_output, serial)
		status := int64(0)
		if api.UserCancel(job.JobID) == true {
			status = STATUS_USERCANCELLED
		} else {
			status = STATUS_ERROR
		}
		if err := api.sendStatus(job, JobOut{
			Status: status,
			StatusReason: fmt.Sprintf("Script, '%s', exited with error status ('%s')",
				job.ScriptName, err.Error()),
			StatusPercent: 0,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: (Script: '%s') %s", job.ScriptName,
				err.Error()))
		}
		api.RemoveJob(job.JobID)
		return
	}

	status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if status == 0 {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_OK,
			StatusReason:  "Script finished successfully",
			StatusPercent: 100,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
	} else {
		if err := api.sendStatus(job, JobOut{
			Status:        STATUS_ERROR,
			StatusReason:  "Non-zero exit status. Check the log.",
			StatusPercent: 100,
			Errors:        0,
		}); err != nil {
			logit(fmt.Sprintf("Error: %s", err.Error()))
		}
	}

	api.RemoveJob(int64(job.JobID))

	// logout
}

/*
func main() {
    exec_cmd ( os.Args[1:]... )
}
*/
