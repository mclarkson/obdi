// Obdi - a REST interface and GUI for deploying software
// Copyright (C) 2014  Mark Clarkson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
			r := regexp.MustCompile(`(?s)([^ ]+=)["'](.*)["']`)
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
		// A user job (the default, should be type=1)
		line, err = rdr.ReadString('\n')
		api.sendOutputLine(job, line, serial)
		for err == nil {
			serial++
			a := ""
			starttime := time.Now()
			// Up to 50 lines can be read in one hit
			for i := 0; i < 50; i++ {
				if time.Duration(time.Since(starttime)) > time.Second {
					break
				}
				line, err = rdr.ReadString('\n')
				a = a + line
			}
			api.sendOutputLine(job, a, serial)
		}
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
