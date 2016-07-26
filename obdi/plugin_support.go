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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func (api *Api) RunPluginUsingRPC(w rest.ResponseWriter, r *rest.Request,
	pluginFile, port, queryType string) ([]byte, error) {

	// Exec the go plugin
	// Look for ./plugins/<endpoint>/<subitem>
	// or ./plugins/<endpoint>

	reply := []byte{}

	cmd := exec.Command(pluginFile, port)

	// Set up buffer for stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		txt := fmt.Sprintf("exec.StdoutPipe error. %s", err)
		logit(txt)
		rest.Error(w, txt, 400)
		return reply, ApiError{"Error"}
	}
	rdr := bufio.NewReader(stdout)

	err = cmd.Start()
	if err != nil {
		txt := fmt.Sprintf("exec.Command error. %s", err)
		logit(txt)
		rest.Error(w, txt, 400)
		return reply, ApiError{"Error"}
	}

	// Secret Data Token - used by plugins to access jsonobjects.
	// Create a GUID and add it to sessions for this user.

	user := User{}
	sdtoken_user := "sduser"
	session := Session{}
	SDToken := ""
	mutex.Lock()
	if api.db.First(&user, "login = ?", sdtoken_user).RecordNotFound() {
		mutex.Unlock()
		logit("'sduser' not found. Access to secret data is NOT possible.")
	} else {
		mutex.Unlock()

		SDToken = NewGUID()

		session = Session{
			Guid:   SDToken,
			UserId: user.Id,
		}
		mutex.Lock()
		if err := api.db.Save(&session).Error; err != nil {
			rest.Error(w, err.Error(), 400)
			mutex.Unlock()
			return reply, ApiError{"Error"}
		}
		mutex.Unlock()
	}

	// Make the RPC call

	type Args struct {
		PathParams  map[string]string
		QueryString map[string][]string
		PostData    []byte
		QueryType   string
		SDToken     string // secret data token
	}

	// Give it 5ms to start
	time.Sleep(3 * time.Millisecond)

	client, err := rpc.Dial("tcp", ":"+port)
	numtries := 1
	for {
		if err != nil {
			// Retry every millisecond
			time.Sleep(5 * time.Millisecond)
			//logit( "Retrying connection to port " + port )
			client, err = rpc.Dial("tcp", ":"+port)
			if numtries > 2000 {
				txt := fmt.Sprintf("Could not connect to plugin"+
					" on port "+port, err)
				rest.Error(w, txt, 400)
				logit(txt)
				// TODO: Kill the client!
				err = cmd.Wait()
				return reply, ApiError{"Error"}
			}
			numtries += 1
			continue
		}
		break
	}

	var postData []byte

	switch queryType {
	case "GET": //do nowt
	case "POST":
		// We don't know how to decode the data so save it raw
		postData, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			txt := fmt.Sprintf("Could not get post data. %s", err)
			rest.Error(w, txt, 400)
			logit(txt)
			err = cmd.Wait()
			return reply, ApiError{"Error"}
		}
	case "PUT":
		// We don't know how to decode the data so save it raw
		postData, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			txt := fmt.Sprintf("Could not get post data. %s", err)
			rest.Error(w, txt, 400)
			logit(txt)
			err = cmd.Wait()
			return reply, ApiError{"Error"}
		}
	case "DELETE": //do nowt
	}

	args := &Args{r.PathParams, r.URL.Query(), postData, queryType, SDToken}
	args.PathParams["PluginDatabasePath"] = config.PluginDbPath
	err = client.Call("Plugin.HandleRequest", args, &reply)
	if err != nil {
		txt := fmt.Sprintf("client.Call error. %s", err)
		rest.Error(w, txt, 400)
		logit(txt)
		err = cmd.Wait()
		return reply, ApiError{"Error"}
	}
	client.Close()

	line := ""
	for err == nil {
		line, err = rdr.ReadString('\n')
		if len(line) > 2 {
			logit(line)
		}
	}

	err = cmd.Wait()

	// We delete the sdtoken after the plugin exits
	mutex.Lock()
	if err := api.db.Delete(&session).Error; err != nil {
		rest.Error(w, err.Error(), 400)
		mutex.Unlock()
		return reply, ApiError{"Error"}
	}
	mutex.Unlock()

	return reply, nil
}

func (api *Api) CompilePlugin(w rest.ResponseWriter,
	pluginFile, endpoint, subitem string) error {

	api.compile.Lock()

	sourceDir := path.Join(config.GoPluginSource, endpoint, "go")
	sourceFile := path.Join(sourceDir, subitem+".go")
	clientLib := path.Join(sourceDir, "obdi_clientlib.go")

	compile := false

	if statExe, err := os.Stat(pluginFile); os.IsNotExist(err) {
		compile = true
	} else {
		if statSrc, err := os.Stat(sourceFile); !os.IsNotExist(err) {
			if statExe.ModTime().Before(statSrc.ModTime()) {
				logit("Plugin is older than source. Recompiling plugin.")
				compile = true
			} else {
				compile = false
			}
		}
	}

	if compile == true {
		if _, err := os.Stat(clientLib); os.IsNotExist(err) {
			clientLib = ""
		}

		logit("Plugin does not exist. Compiling plugin.")
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			txt := fmt.Sprintf("Plugin endpoint '%s/%s' does not exist.",
				endpoint, subitem)
			logit(txt)
			rest.Error(w, txt, 400)
			api.compile.Unlock()
			return ApiError{"Error"}
		}
		// Make plugin directory
		if err := os.Mkdir(path.Dir(pluginFile), 01750); err != nil {
			if !os.IsExist(err) {
				txt := fmt.Sprintf("Plugin endpoint '%s/%s' does not exist. "+
					"Compile failed creating the directory"+
					" '%s'. System said '%s'.",
					endpoint, subitem, pluginFile, err)
				logit(txt)
				rest.Error(w, txt, 400)
				api.compile.Unlock()
				return ApiError{"Error"}
			}
		}

		// Compile

		os.Setenv("PATH", "/usr/bin:/bin:"+path.Join(config.GoRoot, "bin"))
		var cmd *exec.Cmd
		if len(clientLib) > 0 {
			cmd = exec.Command("go", "build", "-o", pluginFile, sourceFile, clientLib)
		} else {
			cmd = exec.Command("go", "build", "-o", pluginFile, sourceFile)
		}
		// Need GOROOT and GOPATH to compile
		cmd.Env = append(cmd.Env, "GOROOT="+config.GoRoot)
		cmd.Env = append(cmd.Env, "GOPATH="+sourceDir)
		// Can't find the go binary without PATH being set
		cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))
		var sout bytes.Buffer
		var serr bytes.Buffer
		cmd.Stdout = &sout
		cmd.Stderr = &serr

		err := cmd.Run()
		if err != nil {
			txt := fmt.Sprintf("Plugin endpoint '%s/%s' does not exist. "+
				"Compile failed for '%s'."+
				" System said '%s'. STDOUT: %s STDERR: %s",
				endpoint, subitem, pluginFile, err,
				cmd.Stdout, cmd.Stderr)
			logit(txt)
			rest.Error(w, txt, 400)
			api.compile.Unlock()
			return ApiError{"Error"}
		}
	}

	api.compile.Unlock()

	return nil
}

func (api *Api) GenericGetEndpoint(w rest.ResponseWriter, r *rest.Request) {

	// Reserve the tcp plugin port now
	iport := api.Port()
	port := strconv.FormatInt(iport, 10)
	defer api.DecrementPort(iport)

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	var errl error = nil
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		logit(errl.Error())
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	// If the plugin isn't available try to compile it

	endpoint := r.PathParam("endpoint")
	subitem := r.PathParam("subitem")
	pluginFile := path.Join(config.GoPluginDir, endpoint, subitem)
	if err := api.CompilePlugin(w, pluginFile, endpoint, subitem); err != nil {
		return // Just return full error was written in CompilePlugin
	}

	// Run the Go plugin

	var reply []byte
	var err error

	if reply, err = api.RunPluginUsingRPC(w, r, pluginFile, port, "GET"); err != nil {
		return // Just return. Full error was written in RunPluginUsingRPC
	}

	// Decode arbitrary JSON. Pull out the mandatory PluginReturn
	// and PluginError fields - all plugins send at least these.

	var generic interface{}
	json.Unmarshal(reply, &generic)
	genericReply := generic.(map[string]interface{})

	var pluginReturn int64
	var pluginError string
	for k, v := range genericReply {
		if k == "PluginReturn" {
			pluginReturn = int64(v.(float64))
		}
		if k == "PluginError" {
			pluginError = v.(string)
		}
	}

	// Return, 0 - success, 1 - error
	if pluginReturn == 1 {
		txt := fmt.Sprintf("Plugin returned error. %s", pluginError)
		rest.Error(w, txt, 400)
		logit(txt)
		return
	}

	// Response is already json encoded so send it raw
	w.(http.ResponseWriter).Write(reply)

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
}

func (api *Api) GenericPostEndpoint(w rest.ResponseWriter, r *rest.Request) {

	// Reserve the tcp plugin port now
	iport := api.Port()
	port := strconv.FormatInt(iport, 10)
	defer api.DecrementPort(iport)

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	var errl error = nil
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		logit(errl.Error())
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	// If the plugin isn't available try to compile it

	endpoint := r.PathParam("endpoint")
	subitem := r.PathParam("subitem")
	pluginFile := path.Join(config.GoPluginDir, endpoint, subitem)
	if err := api.CompilePlugin(w, pluginFile, endpoint, subitem); err != nil {
		return
	}

	// Run the Go plugin

	var reply []byte
	var err error

	if reply, err = api.RunPluginUsingRPC(w, r, pluginFile, port, "POST"); err != nil {
		return // Just return full error was written in RunPluginUsingRPC
	}

	// Decode arbitrary JSON. Pull out the mandatory PluginReturn
	// and PluginError fields - all plugins send at least these.

	var generic interface{}
	json.Unmarshal(reply, &generic)
	genericReply := generic.(map[string]interface{})

	var pluginReturn int64
	var pluginError string
	for k, v := range genericReply {
		if k == "PluginReturn" {
			pluginReturn = int64(v.(float64))
		}
		if k == "PluginError" {
			pluginError = v.(string)
		}
	}

	// Return, 0 - success, 1 - error
	if pluginReturn == 1 {
		txt := fmt.Sprintf("Plugin returned error. %s", pluginError)
		rest.Error(w, txt, 400)
		logit(txt)
		return
	}

	// Response is already json encoded so send it raw
	w.(http.ResponseWriter).Write(reply)

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
}

func (api *Api) GenericDeleteEndpoint(w rest.ResponseWriter, r *rest.Request) {

	// Reserve the tcp plugin port now
	iport := api.Port()
	port := strconv.FormatInt(iport, 10)
	defer api.DecrementPort(iport)

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	var errl error = nil
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		logit(errl.Error())
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	// If the plugin isn't available try to compile it

	endpoint := r.PathParam("endpoint")
	subitem := r.PathParam("subitem")
	pluginFile := path.Join(config.GoPluginDir, endpoint, subitem)
	if err := api.CompilePlugin(w, pluginFile, endpoint, subitem); err != nil {
		return
	}

	// Run the Go plugin

	var reply []byte
	var err error

	if reply, err = api.RunPluginUsingRPC(w, r, pluginFile, port, "DELETE"); err != nil {
		return // Just return full error was written in RunPluginUsingRPC
	}

	// Decode arbitrary JSON. Pull out the mandatory PluginReturn
	// and PluginError fields - all plugins send at least these.

	var generic interface{}
	json.Unmarshal(reply, &generic)
	genericReply := generic.(map[string]interface{})

	var pluginReturn int64
	var pluginError string
	for k, v := range genericReply {
		if k == "PluginReturn" {
			pluginReturn = int64(v.(float64))
		}
		if k == "PluginError" {
			pluginError = v.(string)
		}
	}

	// Return, 0 - success, 1 - error
	if pluginReturn == 1 {
		txt := fmt.Sprintf("Plugin returned error. %s", pluginError)
		rest.Error(w, txt, 400)
		logit(txt)
		return
	}

	// Response is already json encoded so send it raw
	w.(http.ResponseWriter).Write(reply)

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
}

func (api *Api) GenericPutEndpoint(w rest.ResponseWriter, r *rest.Request) {

	// Reserve the tcp plugin port now
	iport := api.Port()
	port := strconv.FormatInt(iport, 10)
	defer api.DecrementPort(iport)

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	var errl error = nil
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		logit(errl.Error())
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	// If the plugin isn't available try to compile it

	endpoint := r.PathParam("endpoint")
	subitem := r.PathParam("subitem")
	pluginFile := path.Join(config.GoPluginDir, endpoint, subitem)
	if err := api.CompilePlugin(w, pluginFile, endpoint, subitem); err != nil {
		return
	}

	// Run the Go plugin

	var reply []byte
	var err error

	if reply, err = api.RunPluginUsingRPC(w, r, pluginFile, port, "PUT"); err != nil {
		return // Just return full error was written in RunPluginUsingRPC
	}

	// Decode arbitrary JSON. Pull out the mandatory PluginReturn
	// and PluginError fields - all plugins send at least these.

	var generic interface{}
	json.Unmarshal(reply, &generic)
	genericReply := generic.(map[string]interface{})

	var pluginReturn int64
	var pluginError string
	for k, v := range genericReply {
		if k == "PluginReturn" {
			pluginReturn = int64(v.(float64))
		}
		if k == "PluginError" {
			pluginError = v.(string)
		}
	}

	// Return, 0 - success, 1 - error
	if pluginReturn == 1 {
		txt := fmt.Sprintf("Plugin returned error. %s", pluginError)
		rest.Error(w, txt, 400)
		logit(txt)
		return
	}

	// Response is already json encoded so send it raw
	w.(http.ResponseWriter).Write(reply)

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
}
