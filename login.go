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

import (
	//"fmt"
	"github.com/mclarkson/deployman/external/ant0ine/go-json-rest/rest"
)

// DoLogin processes "POST /login" queries.
//
// Checks login name and passhash stored in database.
// If correct then 200 header and GUID are sent.
//    any previous sessions for that user are closed.
//    new session entry is made in session table.
// If not correct then 400 header with error message.
//
func (api *Api) DoLogin(w rest.ResponseWriter, r *rest.Request) {

	// Get the Login and Password

	userData := struct{ Login, Password string }{}

	if err := r.DecodeJsonPayload(&userData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if len(userData.Login) == 0 || len(userData.Password) == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}

	//fmt.Printf( "\n%#v\n", userData )
	// Get passhash for login from database

	user := User{}
	if api.db.Where(User{Login: userData.Login}).
		First(&user).RecordNotFound() {
		rest.Error(w, "User or password error.", 400)
		return
	}

	// Check password against hash

	c := &Crypt{}
	c.Pass = []byte(userData.Password)
	c.Hash = []byte(user.Passhash)
	if err := c.Check(); err != nil {
		rest.Error(w, "User or password error.", 400)
		return
	}

	// The user's password matches.
	// Delete old session(s) and create a new one.

	guid := NewGUID()
	session := Session{}

	for {
		session = Session{}
		if api.db.Where(Session{UserId: user.Id}).
			First(&session).RecordNotFound() {

			session = Session{
				Guid:   guid,
				UserId: user.Id,
			}
			if err := api.db.Save(&session).Error; err != nil {
				rest.Error(w, err.Error(), 400)
				return
			}
			break

		} else {

			if err := api.db.Delete(&session).Error; err != nil {
				rest.Error(w, err.Error(), 400)
				return
			}
		}
	}

	logit("User '" + user.Login + "' logged in")
	api.LogActivity(session.Id, "User '"+user.Login+"' logged in.")

	w.WriteJson(struct{ GUID string }{guid})
}

func (api *Api) Logout(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	session := Session{}
	var errl error = nil
	if session, errl = api.CheckLogin(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 400)
		return
	}

	if err := api.db.Delete(&session).Error; err != nil {
		rest.Error(w, err.Error(), 400)
		return
	}

	logit("User '" + login + "' logged out")
	api.LogActivity(session.Id, "User '"+login+"' logged out")
	w.WriteJson("Success")
}
