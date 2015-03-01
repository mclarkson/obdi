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

// All api calls have the username and GUID to be sent as part of the request

import (
	//"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"os/exec"
	"syscall"
  "path"
  "fmt"
  "os"
  "strings"
  "encoding/json"
  "io/ioutil"
  "sync"
)

var cachemutex = &sync.Mutex{}

type RepoType struct {
  Name              string
  Description       string
  Url               string
  // Multiple dependencies
  Depends           []DependsType
  // Each repo can have multiple version referenced by the commit sha
  Versions          []VersionsType
}

type DependsType struct {
  Name              string
  Version           string
  VersionMatchType  string
}

type VersionsType struct {
  Version           string
  CommitSHA         string
  CodeName          string
  Type              string
  ObdiCompatibility ObdiCompatibilityType
}

type ObdiCompatibilityType struct {
  Version           string
  Type              string
}

func (api *Api) GetAllRepoPlugins(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Anyone can view plugins

	/*
	   if login != "admin" {
	       rest.Error(w, "Not allowed", 400)
	       return
	   }
	*/

	//session := Session{}
	var errl error = nil
	//if session,errl = api.CheckLogin( login, guid ); errl != nil {
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

  // Get the list of remote repositories to check

	repos := []Repo{}
	qs := r.URL.Query() // Query string - map[string][]string

  mutex.Lock()
  dbresult := api.db.Order("id").Find(&repos)
  mutex.Unlock()
  if dbresult.Error != nil {
    if !dbresult.RecordNotFound() {
      rest.Error(w, dbresult.Error.Error(), 500)
      return
    }
		rest.Error(w, "No repositories are configured.", 400)
    return
  }

  // Check for git binary
  _, err := exec.LookPath("git")
  if err != nil {
		rest.Error(w, "Could not find 'git' command. Please install git.", 400)
    return
  }

  // Test the remote repositories

  for i := range repos {
    url := repos[i].Url

    cmd := exec.Command("git", "ls-remote", url )
    //cmd.Dir = path.Join(config.StaticContent, "plugins")
    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Git command failed: git ls-remote " + url, 400)
      return
    }
  }

  // Get the metadata list(s)

  cachemutex.Lock()
  defer cachemutex.Unlock()

  for i := range repos {
    url := repos[i].Url

    // cd
    if err := os.Chdir( config.CacheDir ); err != nil {
      rest.Error(w, "Could not change into directory: " + config.CacheDir, 400)
      return
    }

    // http://address/repodir.git -> repodir
    tmpval := strings.Split( url, "/" )
    repodirname := strings.TrimRight( tmpval[len(tmpval)-1], ".git" )

    // rm -rf repodir
    repodir := path.Join( config.CacheDir,repodirname )
    if err := os.RemoveAll( repodir ); err != nil {
      rest.Error(w, "Could not delete directory: " + repodir, 400)
      return
    }

    cmd := exec.Command("git", "clone", url, repodir)

		cmd.Env = []string{}
    if len(config.HttpProxy) > 0 {
      cmd.Env = append(cmd.Env, "http_proxy="+config.HttpProxy)
    }
    if len(config.HttpsProxy) > 0 {
      cmd.Env = append(cmd.Env, "https_proxy="+config.HttpsProxy)
    }

    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Git command failed: git clone " + url, 400)
      return
    }
  }

  // Read and join the repodata.json files

  /* For example
  {
    "Name":"saltregexmanager",
    "Description":"Plugin to allow Salt state ... expressions.",
    "Url":"https://github.com/mclarkson/obdi-saltregexmanager.git",
    "Depends":[{
      "Name":"salt",
      "Version":"0.1.3",
      "VersionMatchType":"AtLeast"
    }],
    "Versions":[
      {
        "Version":"0.1.3",
        "CommitSHA":"5ae1e2c7496669ca63391d87f04c426f2e7dd6ce",
        "CodeName":"china",
        "Type":"alpha",
        "ObdiCompatibility":{"Version":"0.1.3","Type":"AtLeast"}
      }
    ]
  }
  */

  var repodata []RepoType

  for i := range repos {
    // Work out where the repo (repodir) would have been written
    url := repos[i].Url
    tmpval := strings.Split( url, "/" )
    repodirname := strings.TrimRight( tmpval[len(tmpval)-1], ".git" )
    repodir := path.Join( config.CacheDir,repodirname,"repodata.json" )

    // Read the repodata.json file
    file, err := ioutil.ReadFile( repodir )
    if err != nil {
      rest.Error(w, "Could not read repodata.json file: " + err.Error(), 400)
      return
    }
    //fmt.Printf("%s\n", string(file))

    var tmprepodata []RepoType
    err = json.Unmarshal(file, &tmprepodata)
    if err != nil {
      rest.Error(w, "Error decoding json: " + err.Error(), 400)
      return
    }
    repodata = append(repodata,tmprepodata...)
  }

  if len(qs["installable"]) > 0 {

    // Get list of installed plugins and subtract from list of
    // available plugins

    /*
        TODO: Check - Depends, ObdiCompatibility 
    */

    plugins := []Plugin{}

    mutex.Lock()
    err := api.db.Order("name").Find(&plugins)
    mutex.Unlock()
    if err.Error != nil {
      // No results is not an error
      if !err.RecordNotFound() {
        rest.Error(w, err.Error.Error(), 500)
        return
      }
    }

    var repoInstallable []RepoType

RepoLoop:
    for i := range repodata {
      reponame := repodata[i].Name
      for j := range plugins {
        if reponame == plugins[j].Name {
          continue RepoLoop
        }
      }
      repoInstallable = append( repoInstallable, repodata[i] )
    }

    repodata = repoInstallable
  }

  if len(repodata) == 0 {
    repodata = []RepoType{}
  }

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&repodata)
}

func (api *Api) AddRepoPlugin(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Only admin is allowed

	if login != "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	session := Session{}
	var errl error
	if session, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	// Decode the POST data.

  type RepoPlugin struct {
    Name    string
  }

	pluginData := RepoPlugin{}

	if err := r.DecodeJsonPayload(&pluginData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(pluginData.Name) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	plugin := Plugin{}
	mutex.Lock()
	if !api.db.Find(&plugin, "name = ?", pluginData.Name).
		RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

  // Get the list of remote repositories to check

	repos := []Repo{}

  mutex.Lock()
  dbresult := api.db.Order("id").Find(&repos)
  mutex.Unlock()
  if dbresult.Error != nil {
    if !dbresult.RecordNotFound() {
      rest.Error(w, dbresult.Error.Error(), 500)
      return
    }
		rest.Error(w, "No repositories are configured.", 400)
    return
  }

  // Check for git binary
  _, err := exec.LookPath("git")
  if err != nil {
		rest.Error(w, "Could not find 'git' command. Please install git.", 400)
    return
  }

  // Test the remote repositories

  for i := range repos {
    url := repos[i].Url

    cmd := exec.Command("git", "ls-remote", url )
    //cmd.Dir = path.Join(config.StaticContent, "plugins")
    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Git command failed: git ls-remote " + url, 400)
      return
    }
  }

  // Get the metadata list(s)

  cachemutex.Lock()
  defer cachemutex.Unlock()

  for i := range repos {
    url := repos[i].Url

    // cd
    if err := os.Chdir( config.CacheDir ); err != nil {
      rest.Error(w, "Could not change into directory: " + config.CacheDir, 400)
      return
    }

    // http://address/repodir.git -> repodir
    tmpval := strings.Split( url, "/" )
    repodirname := strings.TrimRight( tmpval[len(tmpval)-1], ".git" )

    // rm -rf repodir
    repodir := path.Join( config.CacheDir,repodirname )
    if err := os.RemoveAll( repodir ); err != nil {
      rest.Error(w, "Could not delete directory: " + repodir, 400)
      return
    }

    cmd := exec.Command("git", "clone", url, repodir)

		cmd.Env = []string{}
    if len(config.HttpProxy) > 0 {
      cmd.Env = append(cmd.Env, "http_proxy="+config.HttpProxy)
    }
    if len(config.HttpsProxy) > 0 {
      cmd.Env = append(cmd.Env, "https_proxy="+config.HttpsProxy)
    }

    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Git command failed: git clone " + url, 400)
      return
    }
  }

  // Read and join the repodata.json files

  /* For example
  {
    "Name":"saltregexmanager",
    "Description":"Plugin to allow Salt state ... expressions.",
    "Url":"https://github.com/mclarkson/obdi-saltregexmanager.git",
    "Depends":[{
      "Name":"salt",
      "Version":"0.1.3",
      "VersionMatchType":"AtLeast"
    }],
    "Versions":[
      {
        "Version":"0.1.3",
        "CommitSHA":"5ae1e2c7496669ca63391d87f04c426f2e7dd6ce",
        "CodeName":"china",
        "Type":"alpha",
        "ObdiCompatibility":{"Version":"0.1.3","Type":"AtLeast"}
      }
    ]
  }
  */

  var repodata []RepoType

  for i := range repos {
    // Work out where the repo (repodir) would have been written
    url := repos[i].Url
    tmpval := strings.Split( url, "/" )
    repodirname := strings.TrimRight( tmpval[len(tmpval)-1], ".git" )
    repodir := path.Join( config.CacheDir,repodirname,"repodata.json" )

    // Read the repodata.json file
    file, err := ioutil.ReadFile( repodir )
    if err != nil {
      rest.Error(w, "Could not read repodata.json file: " + err.Error(), 400)
      return
    }
    //fmt.Printf("%s\n", string(file))

    var tmprepodata []RepoType
    err = json.Unmarshal(file, &tmprepodata)
    if err != nil {
      rest.Error(w, "Error decoding json: " + err.Error(), 400)
      return
    }
    repodata = append(repodata,tmprepodata...)
  }

  {

    // Get list of installed plugins and subtract from list of
    // available plugins

    /*
        TODO: Check - Depends, ObdiCompatibility 
    */

    plugins := []Plugin{}

    mutex.Lock()
    err := api.db.Order("name").Find(&plugins)
    mutex.Unlock()
    if err.Error != nil {
      // No results is not an error
      if !err.RecordNotFound() {
        rest.Error(w, err.Error.Error(), 500)
        return
      }
    }

    var repoInstallable []RepoType

RepoLoop:
    for i := range repodata {
      reponame := repodata[i].Name
      for j := range plugins {
        if reponame == plugins[j].Name {
          continue RepoLoop
        }
      }
      repoInstallable = append( repoInstallable, repodata[i] )
    }

    repodata = repoInstallable
  }

  if len(repodata) == 0 {
    repodata = []RepoType{}
  }

  // repodata containes the list of available plugins
  // See if the plugin we want is in the list

  pluginIndex := 0
  pluginFound := false
  for i := range repodata {
    if pluginData.Name == repodata[i].Name {
      pluginFound = true
      pluginIndex = i
    }
  }

  if pluginFound == false {
    rest.Error(w, "Plugin, '" + pluginData.Name + "' not found.", 400)
    return
  }

  //
  // Finally! Get and install the plugin
  //

  {
    url := repodata[pluginIndex].Url

    // cd
    if err := os.Chdir( config.GoPluginSource ); err != nil {
      rest.Error(w, "Could not change into directory: " +
        config.GoPluginSource, 400)
      return
    }

    // rm -rf repodir
    repodir := path.Join( config.GoPluginSource,pluginData.Name )
    if err := os.RemoveAll( repodir ); err != nil {
      rest.Error(w, "Could not delete directory: " + repodir, 400)
      return
    }

    cmd := exec.Command("git", "clone", url, pluginData.Name)

    cmd.Env = []string{}
    if len(config.HttpProxy) > 0 {
      cmd.Env = append(cmd.Env, "http_proxy="+config.HttpProxy)
    }
    if len(config.HttpsProxy) > 0 {
      cmd.Env = append(cmd.Env, "https_proxy="+config.HttpsProxy)
    }

    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Git command failed: git clone " + url, 400)
      return
    }
  }

  // Run the installer script sending 'guid' environment variable

  {
    plugindir := path.Join( config.GoPluginSource,pluginData.Name )

    // cd
    if err := os.Chdir( plugindir ); err != nil {
      rest.Error(w, "Could not change into directory: " +
        config.GoPluginSource, 400)
      return
    }

    cmd := exec.Command("/bin/bash", "install_plugin.sh")

    cmd.Env = []string{}
    cmd.Env = append(cmd.Env, "guid="+guid)

    err := cmd.Start()
    if err != nil {
      txt := fmt.Sprintf("exec.Command error. %s", err)
      logit(txt)
      rest.Error(w, txt, 400)
      return
    }
    err = cmd.Wait()
    status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if status != 0 {
      rest.Error(w, "Command, install_plugin.sh failed.", 400)
      return
    }
  }

	api.LogActivity(session.Id, "Downloaded new plugin '"+pluginData.Name+"'.")
	w.WriteJson(&repodata)
}

func (api *Api) UpdateRepoPlugin(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	//guid := r.PathParam("GUID")

	// Only admin is allowed

	if login != "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	//api.LogActivity(session.Id,
	//	"Updated plugin details for '"+plugin.Name+"'.")

	w.WriteJson("Unimplemented")
}

func (api *Api) DeleteRepoPlugin(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	//guid := r.PathParam("GUID")

	// Only admin is allowed

	if login != "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	//api.LogActivity(session.Id,
	//	"Deleted plugin '"+plugin.Name+"'.")

	w.WriteJson("Unimplemented")
}
