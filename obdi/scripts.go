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
	"fmt"
	//"bytes"
	//"net/url"
	//"time"
	//"github.com/jinzhu/gorm"
	//"database/sql"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

func (api *Api) GetAllScripts(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Anyone can read the list of scripts

	/*
	   if login != "admin" {
	       rest.Error(w, "Not allowed", 400)
	       return
	   }*/

	//session := Session{}
	var errl error = nil
	//if session,errl = api.CheckLogin( login, guid ); errl != nil {
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	scripts := []Script{}
	qs := r.URL.Query() // Query string - map[string][]string
	if len(qs["id"]) > 0 {
		srch := qs["id"][0]
		mutex.Lock()
		api.db.Order("id").Find(&scripts, "id = ?", srch)
		mutex.Unlock()
		/*
		   if api.db.Order("id").
		      Find(&scripts, "id = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else if len(qs["name"]) > 0 {
		srch := qs["name"][0]
		mutex.Lock()
		api.db.Order("name").Find(&scripts, "name = ?", srch)
		mutex.Unlock()
		/*
		   if api.db.Order("name").
		      Find(&scripts, "name = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else {
		// No results is not an error
		mutex.Lock()
		err := api.db.Order("id").Find(&scripts)
		mutex.Unlock()
		if err.Error != nil {
			if !err.RecordNotFound() {
				rest.Error(w, err.Error.Error(), 500)
				return
			}
		}
	}

	// Create a slice of maps from users struct
	// to selectively copy database fields for display

	u := make([]map[string]interface{}, len(scripts))
	for i := range scripts {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = scripts[i].Id
		u[i]["Name"] = scripts[i].Name
		u[i]["Desc"] = scripts[i].Desc
		// 'Source' doesn't go through Unmarshall so
		// is output as base64, good. Use nosource to
		// exclude this field.
		if len(qs["nosource"]) == 0 {
			u[i]["Source"] = scripts[i].Source
		}
		u[i]["Type"] = scripts[i].Type
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddScript(w rest.ResponseWriter, r *rest.Request) {

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

	// Can't add if it exists already

	scriptData := Script{}

	if err := r.DecodeJsonPayload(&scriptData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(scriptData.Source) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	script := Script{}
	mutex.Lock()
	if !api.db.Find(&script, "name = ?", scriptData.Name).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Work out type:
	//   Write to disk then use unix 'file -b' (brief)

	if err := ioutil.WriteFile(os.TempDir()+"/obdi_scriptcheck",
		scriptData.Source, 0644); err != nil {
		scriptData.Type = "Unknown type of script"
	} else {
		runCmd := exec.Command("file", "-b", os.TempDir()+"/obdi_scriptcheck")
		output, err := runCmd.Output()
		if err != nil {
			scriptData.Type = "Unknown type of script"
		} else {
			scriptData.Type = string(output)
		}
	}

	// Add script

	mutex.Lock()
	if err := api.db.Save(&scriptData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	// Try to start the script

	text := fmt.Sprintf("Added new script, %s.", scriptData.Name)
	api.LogActivity(session.Id, text)

	scriptData.Source = []byte{}
	w.WriteJson(scriptData)
}

func (api *Api) UpdateScript(w rest.ResponseWriter, r *rest.Request) {

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

	// Ensure user exists

	id := r.PathParam("id")

	// Check that the id string is a number
	if _, err := strconv.Atoi(id); err != nil {
		rest.Error(w, "Invalid id.", 400)
		return
	}

	// Load data from db, then ...
	script := Script{}
	mutex.Lock()
	if api.db.Find(&script, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&script); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	script_srch := Script{}
	mutex.Lock()
	if !api.db.Find(&script_srch, "name = ? and id != ?",
		script.Name, script.Id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Work out type:
	//   Write to disk then use unix 'file -b' (brief)

	if len(script.Source) > 0 {
		if err := ioutil.WriteFile(os.TempDir()+"/obdi_scriptcheck",
			script.Source, 0644); err != nil {
			script.Type = "Write file failed. Type of script unknown. (" +
				err.Error() + ")"
		} else {
			runCmd := exec.Command("file", "-b",
				os.TempDir()+"/obdi_scriptcheck")
			output, err := runCmd.Output()
			if err != nil {
				script.Type = "Unix 'file' failed. Type of script unknown." +
					" (" + err.Error() + ")"
			} else {
				script.Type = string(output)
			}
		}
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	script.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&script).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Updated data centre details for '"+script.Name+"'.")

	w.WriteJson(script)
}

func (api *Api) DeleteScript(w rest.ResponseWriter, r *rest.Request) {

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

	// Delete

	id := 0
	if id, errl = strconv.Atoi(r.PathParam("id")); errl != nil {
		rest.Error(w, "Invalid id.", 400)
		return
	}

	script := Script{}
	mutex.Lock()
	if api.db.First(&script, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&script).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Deleted data centre '"+script.Name+"'.")

	w.WriteJson("Success")
}
