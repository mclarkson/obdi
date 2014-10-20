// Deployment Manager - a REST interface and GUI for deploying software
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
	"github.com/mclarkson/deployman/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllFiles(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Anyone can view files

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

	files := []File{}
	qs := r.URL.Query() // Query string - map[string][]string
	if len(qs["name"]) > 0 && len(qs["plugin_id"]) > 0 {
		srch1 := qs["name"][0]
		srch2 := qs["plugin_id"][0]
		mutex.Lock()
		api.db.Order("name").Find(&files, "name = ? and plugin_id = ?",
			srch1, srch2)
		mutex.Unlock()
		/*
		   if api.db.Order("name").
		      Find(&files, "name = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else if len(qs["name"]) > 0 {
		srch := qs["name"][0]
		mutex.Lock()
		api.db.Order("name").Find(&files, "name = ?", srch)
		mutex.Unlock()
		/*
		   if api.db.Order("name").
		      Find(&files, "name = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else {
		// No results is not an error
		mutex.Lock()
		err := api.db.Order("name").Find(&files)
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

	u := make([]map[string]interface{}, len(files))
	for i := range files {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = files[i].Id
		u[i]["Name"] = files[i].Name
		u[i]["Desc"] = files[i].Desc
		u[i]["Url"] = files[i].Url
		u[i]["Type"] = files[i].Type
		u[i]["PluginId"] = files[i].PluginId

		plugin := Plugin{}
		mutex.Lock()
		api.db.Model(&files[i]).Related(&plugin)
		mutex.Unlock()

		u[i]["PluginName"] = plugin.Name
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddFile(w rest.ResponseWriter, r *rest.Request) {

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

	fileData := File{}

	if err := r.DecodeJsonPayload(&fileData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(fileData.Name) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	file := File{}
	mutex.Lock()
	if !api.db.Find(&file, "name = ? and plugin_id = ?", fileData.Name,
		fileData.PluginId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Check PluginID

	plugin := Plugin{}
	mutex.Lock()
	if api.db.Find(&plugin, fileData.PluginId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Plugin ID not found", 400)
		return
	}
	mutex.Unlock()

	// Add file

	mutex.Lock()
	if err := api.db.Save(&fileData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id, "Added new file '"+fileData.Name+"'.")
	w.WriteJson(fileData)
}

func (api *Api) UpdateFile(w rest.ResponseWriter, r *rest.Request) {

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
	file := File{}
	mutex.Lock()
	if api.db.Find(&file, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&file); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	file.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&file).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Updated file details for '"+file.Name+"'.")

	w.WriteJson("Success")
}

func (api *Api) DeleteFile(w rest.ResponseWriter, r *rest.Request) {

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

	file := File{}
	mutex.Lock()
	if api.db.First(&file, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&file).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Deleted file '"+file.Name+"'.")

	w.WriteJson("Success")
}
