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

func (api *Api) GetAllEnvs(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Anyone can view envs

	//session := Session{}
	var errl error = nil
	//if session,errl = api.CheckLogin( login, guid ); errl != nil {
	if _, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	defer api.TouchSession(guid)

	envs := []Env{}
	qs := r.URL.Query() // Query string - map[string][]string
	if login == "admin" {
		if len(qs["sys_name"]) > 0 {
			srch := qs["sys_name"][0]
			if len(qs["dc_id"]) > 0 {
				dcid := qs["dc_id"][0]
				mutex.Lock()
				api.db.Order("sys_name").Find(&envs,
					"sys_name = ? and dc_id = ?", srch, dcid)
				mutex.Unlock()
			} else {
				mutex.Lock()
				api.db.Order("sys_name").Find(&envs, "sys_name = ?", srch)
				mutex.Unlock()
			}
			/*
			   if api.db.Order("sys_name").
			      Find(&envs, "sys_name = ?", srch).RecordNotFound() {
			       rest.Error(w, "No results.", 400)
			       return
			   }
			*/
		} else {
			mutex.Lock()
			err := api.db.Order("dc_id,sys_name").Find(&envs)
			mutex.Unlock()
			if err.Error != nil {
				if !err.RecordNotFound() {
					rest.Error(w, err.Error.Error(), 500)
					return
				}
			}
		}
	} else { //Not admin

		// Only return readable/writeable envs for the current user

		additional_where := ""

		if len(qs["env_id"]) > 0 {
			additional_where = "AND envs.id = " + qs["env_id"][0]
		}

		if len(qs["writeable"]) > 0 { // only writeable envs
			mutex.Lock()
			api.db.Where("envs.id in (SELECT perms.env_id from perms "+
				"LEFT JOIN users on users.id=perms.user_id "+
				"WHERE users.login=? and perms.writeable=1) "+
				additional_where,
				login).Find(&envs)
			mutex.Unlock()
		} else { // readable or writeable envs
			mutex.Lock()
			api.db.Where("envs.id in (SELECT perms.env_id from perms "+
				"LEFT JOIN users on users.id=perms.user_id "+
				"WHERE users.login=? and perms.enabled=1) "+
				additional_where,
				login).Find(&envs)
			mutex.Unlock()
		}
	}

	// Create a slice of maps from users struct
	// to selectively copy database fields for display

	u := make([]map[string]interface{}, len(envs))
	for i := range envs {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = envs[i].Id
		u[i]["DispName"] = envs[i].DispName
		u[i]["SysName"] = envs[i].SysName
		//u[i]["WorkerIp"] = envs[i].WorkerIp
		//u[i]["WorkerPort"] = envs[i].WorkerPort
		u[i]["WorkerUrl"] = envs[i].WorkerUrl
		if login == "admin" {
			u[i]["WorkerKey"] = envs[i].WorkerKey
		}
		u[i]["CreatedAt"] = envs[i].CreatedAt

		dc := Dc{}
		mutex.Lock()
		api.db.Model(&envs[i]).Related(&dc)
		mutex.Unlock()

		u[i]["DcSysName"] = dc.SysName
		u[i]["DcDispName"] = dc.DispName
		u[i]["DcId"] = dc.Id
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddEnv(w rest.ResponseWriter, r *rest.Request) {

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

	envData := Env{}

	if err := r.DecodeJsonPayload(&envData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(envData.SysName) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	env := Env{}
	mutex.Lock()
	if !api.db.Find(&env, "sys_name = ? and dc_id = ?",
		envData.SysName, envData.DcId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Add env

	mutex.Lock()
	if err := api.db.Save(&envData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	dc := Dc{}
	mutex.Lock()
	api.db.First(&dc, envData.DcId)
	mutex.Unlock()

	text := fmt.Sprintf("Added new environment '%s->%s'.",
		dc.SysName, envData.SysName)

	api.LogActivity(session.Id, text)
	w.WriteJson(envData)
}

func (api *Api) UpdateEnv(w rest.ResponseWriter, r *rest.Request) {

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
	env := Env{}
	mutex.Lock()
	if api.db.Find(&env, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&env); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	env.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&env).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	dc := Dc{}
	mutex.Lock()
	api.db.First(&dc, env.DcId)
	mutex.Unlock()

	text := fmt.Sprintf("Updated environment '%s->%s'.", dc.SysName, env.SysName)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}

func (api *Api) DeleteEnv(w rest.ResponseWriter, r *rest.Request) {

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

	env := Env{}
	mutex.Lock()
	if api.db.First(&env, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&env).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	dc := Dc{}
	mutex.Lock()
	api.db.First(&dc, env.DcId)
	mutex.Unlock()

	text := fmt.Sprintf("Deleted environment '%s->%s'.",
		dc.SysName, env.SysName)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}
