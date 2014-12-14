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

import (
	"fmt"
	"github.com/mclarkson/obdi/external/jinzhu/gorm"
	"html/template"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	//"strconv"
)

type Api struct {
	db       *gorm.DB
	port     int64
	apimutex *sync.Mutex
	compile  *sync.Mutex
}

type ApiError struct {
	details string
}

// SetDB: Allows to set the gorm.DB
//
func (api *Api) SetDB(db *gorm.DB) {
	api.db = db
}

// Port: Return a port to connect to for RPC and increment it for the
// next connection
func (api *Api) Port() int64 {
	apimutex.Lock()
	portnum := api.port
	api.port += 1
	apimutex.Unlock()
	return portnum
}

func (api *Api) SetPort(portnum int64) {
	apimutex.Lock()
	api.port = portnum
	apimutex.Unlock()
}

func (api *Api) DecrementPort() {
	apimutex.Lock()
	//startport,_ := strconv.ParseInt(config.GoPluginPortStart,10,64)
	if api.port != config.GoPluginPortStart {
		api.port -= 1
	}
	apimutex.Unlock()
}

// TouchSession: Updates the UpdatedAt field
//
// Takes the GUID is a parameter.
//
func (api *Api) TouchSession(guid string) {
	session := Session{}
	mutex.Lock()
	api.db.Where("guid = ?", guid).First(&session)
	session.UpdatedAt = time.Now()
	api.db.Save(&session)
	mutex.Unlock()
}

// LogActivity: Write a message to the activity log
//
func (api *Api) LogActivity(sesId int64, message string) {
	activity := Activity{
		Session_id: sesId,
		Message:    message,
	}
	mutex.Lock()
	api.db.Save(&activity)
	mutex.Unlock()
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%s", e.details)
}

func (api *Api) CheckLoginNoExpiry(login, guid string) (Session, error) {

	user := User{}
	session := Session{}

	// select * from users where login = login
	mutex.Lock()
	if err := api.db.Where("login = ?", login).First(&user).Error; err != nil {
		mutex.Unlock()
		return session, ApiError{"Invalid credentials."}
	}
	// select * from sessions where user_id = user.userid
	if err := api.db.Model(&user).Related(&session).Error; err != nil {
		mutex.Unlock()
		return session, ApiError{"Not logged in."}
	}
	mutex.Unlock()

	// Check GUID
	if session.Guid != guid {
		return session, ApiError{"Invalid GUID."}
	}

	return session, nil
}

func (api *Api) CheckLogin(login, guid string) (Session, error) {

	user := User{}
	session := Session{}

	// select * from users where login = login
	mutex.Lock()
	if err := api.db.Where("login = ?", login).First(&user).Error; err != nil {
		mutex.Unlock()
		return session, ApiError{"Invalid credentials."}
	}
	mutex.Unlock()

	// select * from sessions where user_id = user.userid
	mutex.Lock()
	if err := api.db.Model(&user).Related(&session).Error; err != nil {
		mutex.Unlock()
		return session, ApiError{"Not logged in."}
	}
	mutex.Unlock()

	// Check GUID
	if session.Guid != guid {
		return session, ApiError{"Invalid GUID."}
	}

	// Check session age
	delta := time.Now().Sub(session.UpdatedAt)
	if delta.Minutes() > float64(config.SessionTimeout) {
		return session, ApiError{"Session expired."}
	}

	return session, nil
}

/*
 * The only purpose of serveRunTemplate is to add items to
 * the <head> block. Specifically to add AngularJS controller
 * files to support plugins.
 */
func (api *Api) serveRunTemplate(w http.ResponseWriter, r *http.Request) {

	/* Refer to:
	   http://www.alexedwards.net/blog/serving-static-sites-with-go
	   and
	   http://www.alexedwards.net/blog/a-recap-of-request-handling
	*/

	defaultScripts := []string{}

	// Split the default scripts for admin or run(user)
	// It's alot of unused scripts otherwise.
	if match, _ := path.Match("/manager/admin", r.URL.Path); match == true {
		defaultScripts = []string{
			`js/controllers/login.js`,
			`js/controllers/admin.js`,
			`js/controllers/users.js`,
			`js/controllers/dcs.js`,
			`js/controllers/envs.js`,
			`js/controllers/dccaps.js`,
			`js/controllers/envcaps.js`,
			`js/controllers/scripts.js`,
			`js/controllers/plugins.js`,
		}
	} else {
		// It's /manager/run
		defaultScripts = []string{
			`js/controllers/login.js`,
			`js/controllers/run.js`,
			`js/controllers/sidebar.js`,
		}
	}

	type IndexPageVars struct {
		Items []string
	}

	fp := path.Join(config.StaticContent + "/templates/main-index.html")

	logit(fp)
	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// AngularJS uses {{ }} so we'll use [(
	// The following line creates an unnamed template and sets the
	// Delims for it. Parse files creates a further template named
	// as the file name - this name is used to execute the template
	// and it inherits the Delims. There's no other way to set the
	// Delims!
	tmpl, err := template.New("").Delims(`[(`, `)]`).ParseFiles(fp)
	if err != nil {
		// Log the detailed error
		logit(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Search Files table for controllers that need adding to the
	// HEAD section
	files := []File{}
	mutex.Lock()
	api.db.Order("name").Find(&files, `url != "" and type == 1`)
	mutex.Unlock()

	//fmt.Printf( "%#v\n", files )
	scripts := []string{}
	for i := range files {
		scripts = append(scripts, "plugins/"+files[i].Url)
	}
	defaultScripts = append(defaultScripts, scripts...)

	data := &IndexPageVars{defaultScripts}
	if err := tmpl.ExecuteTemplate(w, path.Base(fp), data); err != nil {
		logit(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func NewApi(db *Database) Api {
	api := Api{}
	api.SetDB(&db.dB)
	//startport,_ := strconv.ParseInt(config.GoPluginPortStart,10,64)
	api.SetPort(config.GoPluginPortStart)
	api.compile = &sync.Mutex{}

	// bigger values allow for more concurrency.
	// increase to avoid 'use of closed network connection' errors
	if config.TransportTimeout < 2 {
		config.TransportTimeout = 2
	}

	return api
}
