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
	"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllJsonObjects(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Only admin is allowed

	if !(login == "admin" || login == "sduser") {
		rest.Error(w, "Not allowed", 400)
		return
	}

	//session := Session{}
	var errl error = nil
	//if session,errl = api.CheckLogin( login, guid ); errl != nil {
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	qs := r.URL.Query() // Query string - map[string][]string

	jsonobjects := []JsonObject{}
	if len(qs["env_id"]) > 0 && len(qs["env_cap_id"]) > 0 {
		mutex.Lock()
		err := api.db.Order("env_id,env_cap_id").Find(&jsonobjects,
			"env_id = ? and env_cap_id = ?", qs["env_id"][0],
			qs["env_cap_id"][0])
		mutex.Unlock()
		if err.Error != nil {
			if !err.RecordNotFound() {
				rest.Error(w, err.Error.Error(), 500)
				return
			}
		}
	} else {
		mutex.Lock()
		err := api.db.Order("env_id,env_cap_id").Find(&jsonobjects)
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

	u := make([]map[string]interface{}, len(jsonobjects))
	for i := range jsonobjects {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = jsonobjects[i].Id
		u[i]["EnvId"] = jsonobjects[i].EnvId
		u[i]["EnvCapId"] = jsonobjects[i].EnvCapId
		u[i]["Json"] = jsonobjects[i].Json

		env_caps := EnvCap{}
		mutex.Lock()
		api.db.Model(&jsonobjects[i]).Related(&env_caps)
		mutex.Unlock()

		u[i]["EnvCapCode"] = env_caps.Code
		u[i]["EnvCapDesc"] = env_caps.Desc
		u[i]["EnvCapId"] = env_caps.Id
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddJsonObject(w rest.ResponseWriter, r *rest.Request) {

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

	jsonobject := JsonObject{}

	if err := r.DecodeJsonPayload(&jsonobject); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(jsonobject.Json) == 0 {
		rest.Error(w, "Incorrect data format received."+
			" Json is unset.", 400)
		return
	}
	{
		tmpjsonobject := JsonObject{}
		mutex.Lock()
		if !api.db.Find(&tmpjsonobject, "env_id = ? and env_cap_id = ?",
			jsonobject.EnvId, jsonobject.EnvCapId).RecordNotFound() {
			mutex.Unlock()
			rest.Error(w, "Record exists.", 400)
			return
		}
		mutex.Unlock()
	}

	// Add jsonobject

	mutex.Lock()
	if err := api.db.Save(&jsonobject).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	env_cap := EnvCap{}
	mutex.Lock()
	api.db.First(&env_cap, jsonobject.EnvCapId)
	mutex.Unlock()

	text := fmt.Sprintf("Added new jsonobject '%s->%s'.",
		env_cap.Desc, jsonobject.Json)

	api.LogActivity(session.Id, text)
	w.WriteJson(jsonobject)
}

func (api *Api) UpdateJsonObject(w rest.ResponseWriter, r *rest.Request) {

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
	jsonobject := JsonObject{}
	mutex.Lock()
	if api.db.Find(&jsonobject, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&jsonobject); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	jsonobject.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&jsonobject).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	env_cap := EnvCap{}
	mutex.Lock()
	api.db.First(&env_cap, jsonobject.EnvCapId)
	mutex.Unlock()

	text := fmt.Sprintf("Updated jsonobject '%s->%s'.",
		env_cap.Desc, jsonobject.Json)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}

func (api *Api) DeleteJsonObject(w rest.ResponseWriter, r *rest.Request) {

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

	jsonobject := JsonObject{}
	mutex.Lock()
	if api.db.First(&jsonobject, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&jsonobject).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	env_cap := EnvCap{}
	mutex.Lock()
	api.db.First(&env_cap, jsonobject.EnvCapId)
	mutex.Unlock()

	text := fmt.Sprintf("Deleted jsonobject '%s->%s'.",
		env_cap.Desc, jsonobject.Json)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}
