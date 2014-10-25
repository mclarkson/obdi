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

func (api *Api) GetAllPerms(w rest.ResponseWriter, r *rest.Request) {

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
	var id int64 = 0

	perms := []Perm{}

	if len(qs["login"]) > 0 {
		user := User{}
		srch := qs["login"][0]
		mutex.Lock()
		err := api.db.First(&user, "login = ?", srch)
		mutex.Unlock()
		if err.Error != nil {
			if !err.RecordNotFound() {
				rest.Error(w, err.Error.Error(), 500)
				return
			}
			if err.RecordNotFound() {
				w.WriteJson(perms)
				return
			}
		}
		id = user.Id
	} else if len(qs["user_id"]) > 0 {
		tid, _ := strconv.Atoi(qs["user_id"][0])
		id = int64(tid)
	}

	if id > 0 {
		mutex.Lock()
		api.db.Order("user_id").Find(&perms, "user_id = ?", id)
		mutex.Unlock()
		/*
		   if api.db.Order("user_id").
		      Find(&perms, "user_id = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else {
		mutex.Lock()
		err := api.db.Order("user_id").Find(&perms)
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

	u := make([]map[string]interface{}, len(perms))
	for i := range perms {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = perms[i].Id
		u[i]["UserId"] = perms[i].UserId
		u[i]["EnvId"] = perms[i].EnvId
		u[i]["Writeable"] = perms[i].Writeable
		u[i]["Enabled"] = perms[i].Enabled
		u[i]["CreatedAt"] = perms[i].CreatedAt

		user := User{}
		mutex.Lock()
		api.db.Model(&perms[i]).Related(&user)
		mutex.Unlock()

		u[i]["UserLogin"] = user.Login
		u[i]["UserEnabled"] = user.Enabled

		env := Env{}
		mutex.Lock()
		api.db.Model(&perms[i]).Related(&env)
		mutex.Unlock()

		u[i]["EnvSysName"] = env.SysName
		u[i]["EnvDispName"] = env.DispName

		dc := Dc{}
		mutex.Lock()
		api.db.Model(&env).Related(&dc)
		mutex.Unlock()

		u[i]["DcSysName"] = dc.SysName
		u[i]["DcDispName"] = dc.DispName
		u[i]["DcId"] = dc.Id
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddPerm(w rest.ResponseWriter, r *rest.Request) {

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

	permData := Perm{}

	if err := r.DecodeJsonPayload(&permData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if permData.UserId == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	perm := Perm{}
	mutex.Lock()
	if !api.db.Find(&perm, "env_id = ? and user_id = ?", permData.EnvId,
		permData.UserId).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Add perm

	mutex.Lock()
	if err := api.db.Save(&permData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Added new environment permission. PermID = '%d'.",
		permData.Id)

	api.LogActivity(session.Id, text)
	w.WriteJson(permData)
}

func (api *Api) UpdatePerm(w rest.ResponseWriter, r *rest.Request) {

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
	perm := Perm{}
	mutex.Lock()
	if api.db.Find(&perm, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&perm); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	perm.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&perm).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Updated environment permission. PermID = '%d'.",
		perm.Id)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}

func (api *Api) DeletePerm(w rest.ResponseWriter, r *rest.Request) {

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

	perm := Perm{}
	mutex.Lock()
	if api.db.First(&perm, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&perm).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	text := fmt.Sprintf("Deleted environment permission. PermID = '%d'.",
		perm.Id)

	api.LogActivity(session.Id, text)

	w.WriteJson("Success")
}
