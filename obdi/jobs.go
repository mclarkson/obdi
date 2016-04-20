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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Job status
const (
	STATUS_UNKNOWN = iota
	STATUS_NOTSTARTED
	STATUS_USERCANCELLED
	STATUS_SYSCANCELLED
	STATUS_INPROGRESS
	STATUS_OK
	STATUS_ERROR
)

/*
 * Send HTTP POST request
 */
func POST(jsondata []byte, url, endpoint string) (r *http.Response, e error) {

	buf := bytes.NewBuffer(jsondata)

	client := &http.Client{Transport: tr}

	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	resp := &http.Response{}

	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	req, err := http.NewRequest("POST",
		url+"/api/"+endpoint, buf)
	if err != nil {
		txt := fmt.Sprintf("Could not send REST request ('%s').", err.Error())
		return resp, ApiError{txt}
	}

	req.Header.Add("Content-Type", `application/json`)

	resp, err = client.Do(req)

	if resp.StatusCode != 200 {
		var body []byte
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
			return resp, ApiError{txt}
		} else {
			body = b
		}
		type myErr struct {
			Error string
		}
		errstr := myErr{}
		if err := json.Unmarshal(body, &errstr); err != nil {
			txt := fmt.Sprintf("Error decoding JSON "+
				"returned from worker - (%s). Check the Worker URL.",
				err.Error())
			return resp, ApiError{txt}
		}

		//txt := fmt.Sprintf("%s", resp.StatusCode)
		return resp, ApiError{errstr.Error}
	} else {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	return resp, nil
}

/*
 * Send HTTP DELETE request
 */
func DELETE(jsondata []byte, url, endpoint string) (r *http.Response,
	e error) {

	buf := bytes.NewBuffer(jsondata)

	client := &http.Client{Transport: tr}

	resp := &http.Response{}

	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	req, err := http.NewRequest("DELETE",
		url+"/api/"+endpoint, buf)
	if err != nil {
		txt := fmt.Sprintf("Could not send REST request ('%s').", err.Error())
		return resp, ApiError{txt}
	}

	req.Header.Add("Content-Type", `application/json`)

	resp, err = client.Do(req)

	return resp, nil
}

func (api *Api) GetAllJobs(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Admin is NOT allowed

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

	jobs := []Job{}
	qs := r.URL.Query() // Query string - map[string][]string
	if len(qs["job_id"]) > 0 {
		srch := qs["job_id"][0]
		mutex.Lock()
		api.db.Order("id desc").Find(&jobs, "id = ?", srch)
		mutex.Unlock()
		/*
		   if api.db.Order("id").
		      Find(&jobs, "id = ?", srch).RecordNotFound() {
		       rest.Error(w, "No results.", 400)
		       return
		   }
		*/
	} else {
		// No results is not an error
		mutex.Lock()
		err := api.db.Order("id desc").Limit(200).Find(&jobs)
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

	u := make([]map[string]interface{}, len(jobs))
	for i := range jobs {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = jobs[i].Id
		u[i]["ScriptId"] = jobs[i].ScriptId
		u[i]["UserLogin"] = jobs[i].UserLogin
		u[i]["Args"] = jobs[i].Args
		u[i]["EnvVars"] = jobs[i].EnvVars
		u[i]["Status"] = jobs[i].Status
		u[i]["StatusReason"] = jobs[i].StatusReason
		u[i]["StatusPercent"] = jobs[i].StatusPercent
		u[i]["CreatedAt"] = jobs[i].CreatedAt
		u[i]["UpdatedAt"] = jobs[i].UpdatedAt
		u[i]["Type"] = jobs[i].Type
		//u[i]["WorkerIp"] = jobs[i].WorkerIp
		//u[i]["WorkerPort"] = jobs[i].WorkerPort
		u[i]["EnvId"] = jobs[i].EnvId
		u[i]["WorkerUrl"] = jobs[i].WorkerUrl
		//u[i]["WorkerKey"] = jobs[i].WorkerKey

		env := Env{}
		mutex.Lock()
		api.db.Model(&jobs[i]).Related(&env)
		mutex.Unlock()

		u[i]["EnvSysName"] = env.SysName
		u[i]["EnvDispName"] = env.DispName
		//u[i]["WorkerUrl"] = env.WorkerUrl

		dc := Dc{}
		mutex.Lock()
		api.db.Model(&env).Related(&dc)
		mutex.Unlock()

		u[i]["DcId"] = dc.Id
		u[i]["DcSysName"] = dc.SysName
		u[i]["DcDispName"] = dc.DispName

		script := Script{}
		mutex.Lock()
		api.db.Model(&jobs[i]).Related(&script)
		mutex.Unlock()

		u[i]["ScriptName"] = script.Name
		u[i]["ScriptDesc"] = script.Desc
		u[i]["ScriptType"] = script.Type
	}

	// Too much noise
	//api.LogActivity( session.Id, "Sent list of users" )
	w.WriteJson(&u)
}

func (api *Api) AddJob(w rest.ResponseWriter, r *rest.Request) {

	// This function will return successfully in all cases, after
	// the sanity checks. Errors are saved in to the job here, or by
	// the worker.

	//
	// TODO: Delete jobs (and output_lines) over N days old here
	//

	logit(fmt.Sprintf("Connection from %s", r.RemoteAddr))

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

	// Can't add if it exists already

	jobData := Job{}

	if err := r.DecodeJsonPayload(&jobData); err != nil {
		rest.Error(w, "Invalid data format received.", 400)
		return
	} else if jobData.ScriptId == 0 {
		rest.Error(w, "Incorrect data format received.", 400)
		return
	}

	jobData.UserLogin = login

	// Sanity checks

	if jobData.ScriptId == 0 {
		txt := "Script ID must be specified"
		rest.Error(w, txt, 400)
		return
	}

	if jobData.EnvId == 0 {
		txt := "Environment ID must be specified"
		rest.Error(w, txt, 400)
		return
	}

	// Add job to DB

	saveJob := func() {
		mutex.Lock()
		if err := api.db.Save(&jobData).Error; err != nil {
			mutex.Unlock()
			rest.Error(w, err.Error(), 400)
			return
		}
		mutex.Unlock()
	}

	saveJob()

	// Get the associated environment data

	env := Env{}
	mutex.Lock()
	api.db.Model(&jobData).Related(&env)
	mutex.Unlock()

	// We need WorkerUrl and WorkerKey

	// If there's an entry in the workers table that matches the current env_id
	// and env_cap_desc (and the cap code is in envcaps and the env is mapped
	// in env_cap_maps) then use that instead of the default.

	for {
		envcaps := []EnvCap{}
		envcapmaps := []EnvCapMap{}
		mutex.Lock()
		api.db.Order("code").Find(&envcaps, "code = ?", jobData.EnvCapDesc)
		mutex.Unlock()
		if len(envcaps) > 0 {
			mutex.Lock()
			api.db.Find(&envcapmaps, "env_id = ? and env_cap_id = ?",
				env.Id, envcaps[0].Id)
			mutex.Unlock()
			if len(envcapmaps) > 0 {
				worker := Worker{}
				if !api.db.Find(&worker, "env_id = ? and env_cap_id = ?",
					env.Id, envcaps[0].Id).RecordNotFound() {

					logit("Found custom worker details for " + jobData.EnvCapDesc)

					jobData.WorkerUrl = worker.WorkerUrl
					jobData.WorkerKey = worker.WorkerKey
					break
				}

				logit("Found " + jobData.EnvCapDesc + " but no workers entry exists.")

			} else {
				logit("Capability Map not found for '" + jobData.EnvCapDesc + "'" +
					", using env_id = '" + strconv.Itoa(int(env.Id)) +
					"' and env_cap_id = '" + strconv.Itoa(int(envcaps[0].Id)) + "'.")
			}
		} else {
			logit("Capability not found for '" + jobData.EnvCapDesc + "'")
		}

		logit("Using default WorkerUrl and WorkerKey.")

		if env.WorkerUrl == "" || env.WorkerKey == "" {
			txt := "WorkerUrl or WorkerKey not set for this environment"
			jobData.Status = STATUS_ERROR
			jobData.StatusReason = txt
			saveJob()
			w.WriteJson(jobData)
			return
		}

		jobData.WorkerUrl = env.WorkerUrl
		jobData.WorkerKey = env.WorkerKey

		break
	}

	// Send the job to the worker

	script := Script{}

	saveJob()

	mutex.Lock()
	if err := api.db.Find(&script, jobData.ScriptId); err.Error != nil {
		mutex.Unlock()
		txt := fmt.Sprintf("Script ID %d not found ('%s')", jobData.ScriptId,
			err.Error.Error())
		jobData.Status = STATUS_ERROR
		jobData.StatusReason = txt
		saveJob()
		w.WriteJson(jobData)
		//rest.Error(w, txt, 400)
		return
	}
	mutex.Unlock()

	// Jobsend definition
	type Jobsend struct {
		ScriptSource []byte
		ScriptName   string
		Args         string
		EnvVars      string
		//NotifURL        string
		JobID int64
		Key   string
		Type  int64 // 1 - user job, 2 - system job
	}

	// Jobsend data
	data := Jobsend{
		ScriptSource: script.Source,
		ScriptName:   script.Name,
		JobID:        jobData.Id,
		Key:          jobData.WorkerKey,
		Args:         jobData.Args,
		EnvVars:      jobData.EnvVars,
		Type:         jobData.Type,
	}

	// Encode
	jsondata, err := json.Marshal(data)
	if err != nil {
		txt := fmt.Sprintf("Error sending job to worker, JSON Encode:",
			err.Error())
		jobData.Status = STATUS_ERROR
		jobData.StatusReason = txt
		saveJob()
		w.WriteJson(jobData)
		//rest.Error(w, txt, 400)
		return
	}
	// POST to worker
	resp, err := POST(jsondata, jobData.WorkerUrl, "jobs")
	//fmt.Printf("%+v", jsondata)
	if err != nil {
		txt := "Could not send job to worker. ('" + err.Error() + "')"
		jobData.Status = STATUS_ERROR
		jobData.StatusReason = txt
		saveJob()
		w.WriteJson(jobData)
		//rest.Error(w, txt, 400)
		return
	}
	defer resp.Body.Close()

	text := fmt.Sprintf("Added new job, %d.", jobData.Id)
	api.LogActivity(session.Id, text)
	w.WriteJson(jobData)
}

func (api *Api) UpdateJob(w rest.ResponseWriter, r *rest.Request) {

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

	// Ensure user exists

	id := r.PathParam("id")

	// Check that the id string is a number
	if _, err := strconv.Atoi(id); err != nil {
		rest.Error(w, "Invalid id.", 400)
		return
	}

	// Load data from db, then ...
	job := Job{}
	mutex.Lock()
	if api.db.Find(&job, id).RecordNotFound() {
		mutex.Unlock()
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Job ID not found.", 400)
		return
	}
	mutex.Unlock()

	// ... overwrite any sent fields
	if err := r.DecodeJsonPayload(&job); err != nil {
		//rest.Error(w, err.Error(), 400)
		rest.Error(w, "Invalid data format received.", 400)
		return
	}

	// Force the use of the path id over an id in the payload
	Id, _ := strconv.Atoi(id)
	job.Id = int64(Id)

	mutex.Lock()
	if err := api.db.Save(&job).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id,
		fmt.Sprintf("Updated job details for jobId %d.", job.Id))

	w.WriteJson(job)
}

func (api *Api) DeleteJob(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Only admin is allowed

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

	job := Job{}
	mutex.Lock()
	if api.db.First(&job, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	mutex.Lock()
	if err := api.db.Delete(&job).Error; err != nil {
		mutex.Unlock()
		rest.Error(w, err.Error(), 400)
		return
	}
	mutex.Unlock()

	api.LogActivity(session.Id, fmt.Sprintf("Deleted job %d.", job.Id))

	w.WriteJson(&job)
}

func (api *Api) KillJob(w rest.ResponseWriter, r *rest.Request) {

	// Check credentials

	login := r.PathParam("login")
	guid := r.PathParam("GUID")

	// Only admin is allowed

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

	job := Job{}
	mutex.Lock()
	if api.db.First(&job, id).RecordNotFound() {
		mutex.Unlock()
		rest.Error(w, "Record not found.", 400)
		return
	}
	mutex.Unlock()

	env := Env{}
	mutex.Lock()
	api.db.Model(&job).Related(&env)
	mutex.Unlock()

	// Add job to DB

	saveJob := func() {
		mutex.Lock()
		if err := api.db.Save(&job).Error; err != nil {
			mutex.Unlock()
			rest.Error(w, err.Error(), 400)
			return
		}
		mutex.Unlock()
	}

	// If there's an entry in the workers table that matches the
	// current env_id and env_cap_desc then use that instead of
	// the default.

	for {
		envcaps := []EnvCap{}
		envcapmaps := []EnvCapMap{}
		mutex.Lock()
		api.db.Order("code").Find(&envcaps, "code = ?", job.EnvCapDesc)
		mutex.Unlock()
		if len(envcaps) > 0 {
			mutex.Lock()
			api.db.Find(&envcapmaps, "env_id = ? and env_cap_id = ?",
				env.Id, envcaps[0].Id)
			mutex.Unlock()
			if len(envcapmaps) > 0 {
				worker := Worker{}
				if !api.db.Find(&worker, "env_id = ? and env_cap_id = ?",
					env.Id, envcaps[0].Id).RecordNotFound() {

					logit("Found custom worker details for " + job.EnvCapDesc)

					job.WorkerUrl = worker.WorkerUrl
					job.WorkerKey = worker.WorkerKey
					break
				}

				logit("Found " + job.EnvCapDesc + " but no workers entry exists.")

			} else {
				logit("Capability Map not found for '" + job.EnvCapDesc + "'" +
					", using env_id = '" + strconv.Itoa(int(env.Id)) +
					"' and env_cap_id = '" + strconv.Itoa(int(envcaps[0].Id)) + "'.")
			}
		} else {
			logit("Capability not found for '" + job.EnvCapDesc + "'")
		}

		logit("Using default WorkerUrl and WorkerKey.")

		if env.WorkerUrl == "" || env.WorkerKey == "" {
			txt := "WorkerUrl or WorkerKey not set for this environment"
			job.Status = STATUS_ERROR
			job.StatusReason = txt
			saveJob()
			w.WriteJson(job)
			return
		}

		job.WorkerUrl = env.WorkerUrl
		job.WorkerKey = env.WorkerKey

		break
	}

	type Jobkill struct {
		JobID int64
		Key   string
	}
	data := Jobkill{
		JobID: job.Id,
		Key:   job.WorkerKey,
	}
	// Encode
	jsondata, err := json.Marshal(data)
	if err != nil {
		txt := fmt.Sprintf(
			"Error sending kill command to worker, JSON Encode:",
			err.Error())
		rest.Error(w, txt, 400)
		return
	}
	// POST to worker
	resp, err := DELETE(jsondata, job.WorkerUrl, "jobs")
	if err != nil {
		txt := "Could not send kill command to worker. ('" + err.Error() + "')"
		rest.Error(w, txt, 400)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var body []byte
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
			rest.Error(w, txt, 400)
			return
		} else {
			body = b
		}
		type myErr struct {
			Error string
		}
		errstr := myErr{}
		if err := json.Unmarshal(body, &errstr); err != nil {
			txt := fmt.Sprintf("Error decoding JSON ('%s')"+
				". Check the Worker URL.", err.Error())
			rest.Error(w, txt, 400)
			return
		}

		txt := "Sending Kill failed. Worker said: '" +
			errstr.Error + "'"
		rest.Error(w, txt, 400)
		return
	} else {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	api.LogActivity(session.Id, fmt.Sprintf("Killed job %d.", job.Id))

	w.WriteJson(&job)
}
