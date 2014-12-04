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
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllEnvCaps(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Only admin is allowed

	if login != "admin" {
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

	envcaps := []EnvCap{}

	if len(qs["code"]) > 0 {
		srch := qs["code"][0]
		mutex.Lock()
		api.db.Order("code").Find(&envcaps, "code = ?", srch)
		mutex.Unlock()
	} else {
		mutex.Lock()
		err := api.db.Order("Code").Find(&envcaps)
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

	u := make([]map[string]interface{}, len(envcaps))
	for i := range envcaps {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = envcaps[i].Id
		u[i]["Code"] = envcaps[i].Code
		u[i]["Desc"] = envcaps[i].Desc
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddEnvCap(w rest.ResponseWriter, r *rest.Request) {

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

	EnvCapData := EnvCap{}

	if err := r.DecodeJsonPayload(&EnvCapData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(EnvCapData.Code) == 0 || len(EnvCapData.Desc) == 0 {
		rest.Error(w, "A required field is empty.", 400)
		return
	}
	EnvCap := EnvCap{}
	mutex.Lock()
	if !api.db.Find(&EnvCap, "code = ?", EnvCapData.Code).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Add EnvCap

	mutex.Lock()
	if err := api.db.Save(&EnvCapData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Added new EnvCap, '%s'.",
		EnvCapData.Code)

	api.LogActivity(session.Id, text)
	w.WriteJson(EnvCapData)
}

func (api *Api) UpdateEnvCap(w rest.ResponseWriter, r *rest.Request) {

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
	EnvCap := EnvCap{}
	mutex.Lock()
	if api.db.Find(&EnvCap, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&EnvCap); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	EnvCap.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&EnvCap).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Updated EnvCap, '%s'.",
		EnvCap.Code)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}

func (api *Api) DeleteEnvCap(w rest.ResponseWriter, r *rest.Request) {

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

	EnvCap := EnvCap{}
	mutex.Lock()
	if api.db.First(&EnvCap, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&EnvCap).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Deleted EnvCap, '%s'.",
		EnvCap.Code)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}
