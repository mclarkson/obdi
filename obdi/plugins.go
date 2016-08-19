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

// All api calls have the username and GUID to be sent as part of the request

import (
	//"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllPlugins(w rest.ResponseWriter, r *rest.Request) {

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

	plugins := []Plugin{}
	qs := r.URL.Query() // Query string - map[string][]string
	if len(qs["name"]) > 0 {
		srch := qs["name"][0]
		mutex.Lock()
		api.db.Order("name").Find(&plugins, "name = ?", srch)
		mutex.Unlock()
		/*
		   if api.db.Order("name").
		      Find(&plugins, "name = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else {
		// No results is not an error
		mutex.Lock()
		err := api.db.Order("name").Find(&plugins)
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

	u := make([]map[string]interface{}, len(plugins))
	for i := range plugins {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = plugins[i].Id
		u[i]["Name"] = plugins[i].Name
		u[i]["Desc"] = plugins[i].Desc
		u[i]["Parent"] = plugins[i].Parent
		u[i]["HasView"] = plugins[i].HasView
		u[i]["CreatedAt"] = plugins[i].CreatedAt
		// UpdatedAt  doesn't get updated 'cos we use Save
		// //u[i]["UpdatedAt"] = plugins[i].CreatedAt
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddPlugin(w rest.ResponseWriter, r *rest.Request) {

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

	pluginData := Plugin{}

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

	// Make sure parent exists
	mutex.Lock()
	if len(pluginData.Parent) > 0 && api.db.Find(&plugin, "name = ?",
		pluginData.Parent).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Parent not found.", 400)
		return
	}
	mutex.Unlock()

	// Add plugin

	mutex.Lock()
	if err := api.db.Save(&pluginData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id, "Added new plugin '"+pluginData.Name+"'.")
	w.WriteJson(pluginData)
}

func (api *Api) UpdatePlugin(w rest.ResponseWriter, r *rest.Request) {

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
	plugin := Plugin{}
	mutex.Lock()
	if api.db.Find(&plugin, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&plugin); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	plugin.Id = int64(Id)

	// Make sure parent exists
	pluginSrch := Plugin{}
	mutex.Lock()
	if len(plugin.Parent) > 0 && api.db.Find(&pluginSrch, "name = ?",
		plugin.Parent).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Parent not found.", 400)
		return
	}
	mutex.Unlock()
	if pluginSrch.Id == plugin.Id {
		rest.Error(w, "Cannot be a parent of itself.", 400)
		return
	}

	mutex.Lock()
	if err := api.db.Save(&plugin).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Updated plugin details for '"+plugin.Name+"'.")

	w.WriteJson("Success")
}

func (api *Api) DeletePlugin(w rest.ResponseWriter, r *rest.Request) {

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

	plugin := Plugin{}
	mutex.Lock()
	if api.db.First(&plugin, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&plugin).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	if err := api.db.Where("plugin_id = ?", plugin.Id).
		Delete(File{}).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Deleted plugin '"+plugin.Name+"'.")

	w.WriteJson("Success")
}
