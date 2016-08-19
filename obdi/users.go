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
	//"bytes"
	//"net/url"
	//"time"
	//"github.com/jinzhu/gorm"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
	"strings"
)

// GetAllUsers processes "GET /users" queries.
//
func (api *Api) GetAllUsers(w rest.ResponseWriter, r *rest.Request) {

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

	users := []User{}
	qs := r.URL.Query() // map[string][]string
	if len(qs["login"]) > 0 {
		srch := qs["login"][0]
		mutex.Lock()
		if api.db.Order("login").
			Find(&users, "login = ?", srch).RecordNotFound() {
			mutex.Unlock()
			rest.Error(w, "No results.", 400)
			return
		}
		mutex.Unlock()
	} else {
		mutex.Lock()
		if api.db.Order("login").Find(&users).RecordNotFound() {
			mutex.Unlock()
			rest.Error(w, "Empty Table.", 400)
			return
		}
		mutex.Unlock()
	}

	// Create a slice of maps from users struct
	// to selectively copy database fields for display

	u := make([]map[string]interface{}, len(users))
	for i := range users {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = users[i].Id
		u[i]["Login"] = users[i].Login
		u[i]["Forename"] = users[i].Forename
		u[i]["Surname"] = users[i].Surname
		u[i]["Enabled"] = users[i].Enabled
		u[i]["CreatedAt"] = users[i].CreatedAt
		u[i]["Email"] = users[i].Email
		u[i]["MultiLogin"] = users[i].MultiLogin
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

// AddUser processes "POST /users" queries.
//
func (api *Api) AddUser(w rest.ResponseWriter, r *rest.Request) {

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

	userData := User{}

	if err := r.DecodeJsonPayload(&userData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(userData.Login) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}
	user := User{}
	mutex.Lock()
	if !api.db.Find(&user, "login = ?", userData.Login).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record exists.", 400)
		return
	}
	mutex.Unlock()

	// Add user

	if len(userData.Passhash) == 0 {
		rest.Error(w, "Empty password not allowed.", 400)
		return
	}

	c := &Crypt{}
	c.Pass = []byte(userData.Passhash)
	c.Crypt()
	userData.Passhash = string(c.Hash)

	mutex.Lock()
	if err := api.db.Save(&userData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id, "Added new user '"+userData.Login+"'.")
	w.WriteJson(userData)
}

// UpdateUser processes "PUT /users" queries.
//
func (api *Api) UpdateUser(w rest.ResponseWriter, r *rest.Request) {

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
	user := User{}
	mutex.Lock()
	if api.db.Find(&user, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	// FIXME: DecodeJsonPayload(&somethingelse) then
	// merge with 'user' manually. This will remove
	// the 'password can't begin with $' limitation.

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&user); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Add user

	if !strings.HasPrefix(user.Passhash, "$") {
		c := &Crypt{}
		c.Pass = []byte(user.Passhash)
		c.Crypt()
		user.Passhash = string(c.Hash)
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	user.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&user).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Updated user details for '"+user.Login+"'.")

	w.WriteJson(user)
}

// DeleteUser processes "DELETE /users" queries.
//
func (api *Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {

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

	user := User{}
	mutex.Lock()
	if api.db.First(&user, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&user).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		"Deleted user '"+user.Login+"'.")

	w.WriteJson("Success")
}
