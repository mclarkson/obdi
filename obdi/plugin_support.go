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
	err := cmd.Start()
	if err != nil {
		txt := fmt.Sprintf("exec.Command error. %s", err)
		logit(txt)
		rest.Error(w, txt, 400)
		return reply, ApiError{"Error"}
	}

	// Make the RPC call

	type Args struct {
		PathParams  map[string]string
		QueryString map[string][]string
		PostData    []byte
		QueryType   string
	}

	// Give it 5ms to start
	time.Sleep(5 * time.Millisecond)

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

	args := &Args{r.PathParams, r.URL.Query(), postData, queryType}
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

	err = cmd.Wait()

	return reply, nil
}

func (api *Api) CompilePlugin(w rest.ResponseWriter,
	pluginFile, endpoint, subitem string) error {

	api.compile.Lock()

	if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
		sourceDir := path.Join(config.GoPluginSource, endpoint, "go")
		sourceFile := path.Join(sourceDir, subitem+".go")
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
		cmd := exec.Command("go", "build", "-o", pluginFile, sourceFile)
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
				" System said '%s'. STDOUT: %s STDERR:",
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
	port := strconv.FormatInt(api.Port(), 10)
	defer api.DecrementPort()

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

func (api *Api) GenericPostEndpoint(w rest.ResponseWriter, r *rest.Request) {

	// Reserve the tcp plugin port now
	port := strconv.FormatInt(api.Port(), 10)
	defer api.DecrementPort()

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
	port := strconv.FormatInt(api.Port(), 10)
	defer api.DecrementPort()

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
	port := strconv.FormatInt(api.Port(), 10)
	defer api.DecrementPort()

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
