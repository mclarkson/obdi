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

import (
	//"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
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
	mutex.Lock()
	if api.db.Where(User{Login: userData.Login}).
		First(&user).RecordNotFound() {
		rest.Error(w, "User or password error.", 400)
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	// Check password against hash

	c := &Crypt{}
	c.Pass = []byte(userData.Password)
	c.Hash = []byte(user.Passhash)
	if err := c.Check(); err != nil {
		rest.Error(w, "User or password error.", 400)
		return
	}

	// The user's password matches.
	// Delete old session(s) (if not MultiLogin) and create a new one

	guid := NewGUID()
	session := Session{}

	for {
		session = Session{}
		mutex.Lock()
		if user.MultiLogin {
			session = Session{
				Guid:   guid,
				UserId: user.Id,
			}
			if err := api.db.Save(&session).Error; err != nil {
				rest.Error(w, err.Error(), 400)
				mutex.Unlock()
				return
			}
			mutex.Unlock()
			break
		} else {
			if api.db.Where(Session{UserId: user.Id}).
				First(&session).RecordNotFound() {

				session = Session{
					Guid:   guid,
					UserId: user.Id,
				}
				if err := api.db.Save(&session).Error; err != nil {
					rest.Error(w, err.Error(), 400)
					mutex.Unlock()
					return
				}
				mutex.Unlock()
				break

			} else {

				if err := api.db.Delete(&session).Error; err != nil {
					rest.Error(w, err.Error(), 400)
					mutex.Unlock()
					return
				}
			}
		}
		mutex.Unlock()
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

	mutex.Lock()
	if err := api.db.Delete(&session).Error; err != nil {
		rest.Error(w, err.Error(), 400)
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	logit("User '" + login + "' logged out")
	api.LogActivity(session.Id, "User '"+login+"' logged out")
	w.WriteJson("Success")
}
