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
	//"bytes"
	//"net/url"
	//"time"
	//"github.com/jinzhu/gorm"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllEnvCapMaps(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Anyone can read environment capability maps
	//
	//if login != "admin" {
	//	rest.Error(w, "Not allowed", 400)
	//	return
	//}

	//session := Session{}
	var errl error = nil
	//if session,errl = api.CheckLogin( login, guid ); errl != nil {
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	qs := r.URL.Query() // Query string - map[string][]string

	envcapmaps := []EnvCapMap{}

	if len(qs["env_id"]) > 0 {
		srch := qs["env_id"][0]
		mutex.Lock()
		api.db.Order("env_id").Find(&envcapmaps, "env_id = ?", srch)
		mutex.Unlock()
	} else {
		mutex.Lock()
		err := api.db.Order("env_id").Find(&envcapmaps)
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

	u := make([]map[string]interface{}, len(envcapmaps))
	for i := range envcapmaps {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = envcapmaps[i].Id
		u[i]["EnvId"] = envcapmaps[i].EnvId
		u[i]["EnvCapId"] = envcapmaps[i].EnvCapId

		env := Env{}
		mutex.Lock()
		api.db.Model(&envcapmaps[i]).Related(&env)
		mutex.Unlock()

		u[i]["EnvSysName"] = env.SysName
		u[i]["EnvDispName"] = env.DispName

		envcap := EnvCap{}
		mutex.Lock()
		api.db.Model(&envcapmaps[i]).Related(&envcap)
		mutex.Unlock()

		u[i]["EnvCapCode"] = envcap.Code
		u[i]["EnvCapDesc"] = envcap.Desc
		u[i]["EnvCapIsWorkerDef"] = envcap.IsWorkerDef
		u[i]["EnvCapIsJsonObjectDef"] = envcap.IsJsonObjectDef
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddEnvCapMap(w rest.ResponseWriter, r *rest.Request) {

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

	envCapMapData := EnvCapMap{}

	if err := r.DecodeJsonPayload(&envCapMapData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if envCapMapData.EnvId == 0 || envCapMapData.EnvCapId == 0 {
		rest.Error(w, "A required field is empty.", 400)
		return
	}
	envCapMap := EnvCapMap{}
	mutex.Lock()
	if !api.db.Find(&envCapMap, "env_id = ? and env_cap_id = ?",
		envCapMapData.EnvId, envCapMapData.EnvCapId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Check that EnvId and EnvCapId exist
	env := Env{}
	mutex.Lock()
	if api.db.Find(&env, envCapMapData.EnvId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Invalid environment id.", 400)
		return
	}
	mutex.Unlock()
	envCap := EnvCap{}
	mutex.Lock()
	if api.db.Find(&envCap, envCapMapData.EnvCapId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Invalid environment capability id.", 400)
		return
	}
	mutex.Unlock()

	// Add EnvCapMap

	mutex.Lock()
	if err := api.db.Save(&envCapMapData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Added new EnvCapMap, '%d'.",
		envCapMapData.Id)

	api.LogActivity(session.Id, text)
	w.WriteJson(envCapMapData)
}

func (api *Api) UpdateEnvCapMap(w rest.ResponseWriter, r *rest.Request) {

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
	envCapMap := EnvCapMap{}
	mutex.Lock()
	if api.db.Find(&envCapMap, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&envCapMap); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	envCapMap.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&envCapMap).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Updated EnvCapMap, '%d'.",
		envCapMap.Id)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}

func (api *Api) DeleteEnvCapMap(w rest.ResponseWriter, r *rest.Request) {

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

	envCapMap := EnvCapMap{}
	mutex.Lock()
	if api.db.First(&envCapMap, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&envCapMap).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Deleted EnvCapMap, '%d'.",
		envCapMap.Id)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}
