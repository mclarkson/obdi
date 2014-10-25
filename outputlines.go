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
	"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"strconv"
)

func (api *Api) GetAllOutputLines(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
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

	outputlines := []OutputLine{}
	qs := r.URL.Query() // Query string - map[string][]string
	if len(qs["job_id"]) > 0 {
		srch := qs["job_id"][0]
		if len(qs["top"]) > 0 {
			mutex.Lock()
			api.db.Order("serial").Limit(qs["top"][0]).Find(&outputlines,
				"job_id = ?", srch)
			mutex.Unlock()
		} else if len(qs["bottom"]) > 0 {
			mutex.Lock()
			// TODO last X lines but *in* order
			api.db.Order("serial desc").Limit(qs["bottom"][0]).
				Find(&outputlines, "job_id = ?", srch)
			mutex.Unlock()
		} else {
			mutex.Lock()
			api.db.Order("serial").Find(&outputlines, "job_id = ?", srch)
			mutex.Unlock()
		}
	} else {
		mutex.Lock()
		err := api.db.Order("serial").Find(&outputlines)
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

	u := make([]map[string]interface{}, len(outputlines))
	for i := range outputlines {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = outputlines[i].Id
		u[i]["Serial"] = outputlines[i].Serial
		u[i]["JobId"] = outputlines[i].JobId
		u[i]["Text"] = outputlines[i].Text
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddOutputLine(w rest.ResponseWriter, r *rest.Request) {

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed
	if login == "admin" {
		rest.Error(w, "Not allowed", 400)
		return
	}

	// Check credentials
	//session := Session{}
	var errl error
	if _, errl = api.CheckLoginNoExpiry(login, guid); errl != nil {
		rest.Error(w, errl.Error(), 401)
		return
	}

	//defer api.TouchSession( guid )

	outputLineData := OutputLine{}

	if err := r.DecodeJsonPayload(&outputLineData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if outputLineData.JobId == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}

	// Add OutputLine

	mutex.Lock()
	if err := api.db.Save(&outputLineData).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	//text := ""
	//fmt.Sprintf( text,"%d",outputLineData.JobId )
	//api.LogActivity( session.Id, "Started outputLine logging for job '"+
	//    text+"'." )

	w.WriteJson("Success")
}

/*
func (api *Api) UpdateOutputLine(w rest.ResponseWriter, r *rest.Request) {
    fmt.Println("In UpdateOutputLine")
    w.WriteJson( "In UpdateOutputLine" )
}
*/

func (api *Api) DeleteOutputLine(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is not allowed

	if login == "admin" {
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

	outputline := OutputLine{}
	mutex.Lock()
	if err := api.db.Where("job_id = ?", id).Delete(&outputline).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		fmt.Sprintf("Deleted outputlines for job %d.", id))

	w.WriteJson("Success")
}
