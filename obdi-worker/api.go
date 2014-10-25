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
    "sync"
    "bytes"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
    "time"
)

// Inbound
type JobIn struct {
    ScriptSource    []byte      // From manager
    Args            string      // From manager
    EnvVars         string      // From manager
    //NotifURL        string      // From manager
    JobID           int64       // From manager
    Key             string      // From manager
    Type            int64       // From manager: 1 - user job, 2 - system job
    Guid            string      // Locally created
    Pid             int64       // Locally created
    Errors          int64       // Locally created
    UserCancel      bool        // Used locally only
}

// Outbound: All created locally
type JobOut struct {
    Id              int64       // JobID
    Status          int64
    StatusReason    string
    StatusPercent   int64
    Errors          int64
}

type OutputLine struct {
    Id              int64
    Serial          int64
    JobId           int64
    Text            string
    //Type            int64       // 0 - output, 1 - error output
}

type Api struct {
    jobs        []JobIn
    mutex       *sync.Mutex
    loginmutex  *sync.Mutex
    guid        string
}

type ApiError struct {
    details string
}

const (
    STATUS_UNKNOWN          = iota
    STATUS_NOTSTARTED
    STATUS_USERCANCELLED
    STATUS_SYSCANCELLED
    STATUS_INPROGRESS
    STATUS_OK
    STATUS_ERROR
)

func (api *Api) sendOutputLine( job JobIn, line string, serial int64 ) error {

    data := OutputLine{}
    data.Serial = serial
    data.JobId = job.JobID
    data.Text = line

    jsondata, err := json.Marshal(data)
    if err != nil {
        return ApiError{ "Internal error: Manager login, JSON Encode" }
    }

    _, err = POST(jsondata,
        config.User+"/"+api.Guid()+"/outputlines")
    if err != nil {
        return ApiError{ err.Error() }
    }

    return nil
}

func (api *Api) sendStatus( jobin JobIn, jobout JobOut ) error {

    tries := 0

    r := &http.Response{}

    for {
        jsondata, err := json.Marshal(jobout)
        if err != nil {
            return ApiError{ "Internal error: sendStatus, JSON Encode" }
        }

        resp, err := PUT( jsondata, config.User + "/" + api.Guid() + "/" +
                          "jobs/" + fmt.Sprintf("%d",jobin.JobID) )
        if err != nil {
            return ApiError{ err.Error() }
        }
        r = resp
        // Retry login (only once) on a 401
        if resp.StatusCode != 401 { break }
        if tries == 1 { break }
        resp.Body.Close()
        tries = tries + 1
        api.Login()
    }

    if r.StatusCode != 200 {

        // There was an error
        // Read the response body for details

        var body []byte
        if b, err := ioutil.ReadAll(r.Body); err != nil {
            txt := fmt.Sprintf("Error reading Body ('%s').", err.Error() )
            return ApiError{ txt }
        } else {
            body = b
        }
        type myErr struct {
            Error   string
        }
        errstr := myErr{}
        if err := json.Unmarshal(body,&errstr); err != nil {
            txt := fmt.Sprintf("Error decoding JSON ('%s')", err.Error() )
            return ApiError{ txt }
        }
        txt := fmt.Sprintf("SendStatus to Manager failed ('%s').",
        errstr.Error)
        return ApiError{ txt }

    }

    r.Body.Close()

    return nil
}

func (e ApiError) Error() string {
    return fmt.Sprintf("%s", e.details)
}

type Login struct {
    Login       string
    Password    string
}

func POST(jsondata []byte, endpoint string)(r *http.Response,
e error) {

    buf := bytes.NewBuffer(jsondata)

    // accept bad certs
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    client.Timeout = time.Duration(config.TransportTimeout) * 1e9

    resp, err := client.Post( config.ManUrlPrefix+"/api/"+endpoint,
        "application/json", buf )
    if err != nil {
        txt := fmt.Sprintf("Could not send REST request ('%s').",err.Error())
        return resp, ApiError{ txt }
    }

    return resp, nil
}

func PUT(jsondata []byte, endpoint string)(r *http.Response,
e error) {

    buf := bytes.NewBuffer(jsondata)

    // accept bad certs
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    client.Timeout = time.Duration(config.TransportTimeout) * 1e9

    resp := &http.Response{}

    req, err := http.NewRequest( "PUT",
        config.ManUrlPrefix+"/api/"+endpoint, buf )
    if err != nil {
        txt := fmt.Sprintf("Could not send REST request ('%s').",err.Error())
        return resp, ApiError{ txt }
    }

    req.Header.Add("Content-Type", `application/json`)

    resp, err = client.Do( req )

    return resp, nil
}

func (api *Api) Login() error {

    data := Login{}
    data.Login = config.User
    data.Password = config.Password
    jsondata, err := json.Marshal(data)
    if err != nil {
        return ApiError{ "Internal error: Manager login, JSON Encode" }
    }

    resp, err := POST(jsondata,"login")
    if err != nil {
        logit( "Logging in. POST FAILED." + err.Error() +
               " (Maybe transport_timeout is too low or not set?)" )
        return ApiError{ err.Error() }
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {

        // There was an error
        // Read the response body for details

        var body []byte
        if b, err := ioutil.ReadAll(resp.Body); err != nil {
            txt := fmt.Sprintf("Error reading Body ('%s').", err.Error() )
            return ApiError{ txt }
        } else {
            body = b
        }
        type myErr struct {
            Error   string
        }
        errstr := myErr{}
        if err := json.Unmarshal(body,&errstr); err != nil {
            txt := fmt.Sprintf("Error decoding JSON ('%s')", err.Error() )
            return ApiError{ txt }
        }
        txt := fmt.Sprintf("Login from worker to Manager failed ('%s').",
        errstr.Error)
        return ApiError{ txt }

    } else {

        // Logged in OK, save the GUID

        var body []byte
        if b, err := ioutil.ReadAll(resp.Body); err != nil {
            txt := fmt.Sprintf("Error reading Body ('%s').", err.Error() )
            return ApiError{ txt }
        } else {
            body = b
        }
        type login struct { GUID string }
        loginDetails := login{}
        if err := json.Unmarshal(body,&loginDetails); err != nil {
            txt := fmt.Sprintf(
                "Error decoding login response from Manager ('%s')",
                err.Error() )
            return ApiError{ txt }
        }

        api.UpdateGuid( loginDetails.GUID )
    }

    return nil
}

func (api *Api) Logout( ) error {

    jsondata := []byte{}  // No json for logout

    _, err := POST(jsondata,
        config.User+"/"+api.Guid()+"/logout")
    if err != nil {
        return ApiError{ err.Error() }
    }

    api.UpdateGuid( "" )

    return nil
}

func (api *Api) AppendJob( job JobIn ) {
    api.mutex.Lock()
    api.jobs = append(api.jobs, job)
    api.mutex.Unlock()
}

func (api *Api) FindJob( jobid int64 ) (jobret JobIn, e error ) {
    api.mutex.Lock()
    found := 0
    for i,job := range api.jobs {
        if job.JobID == jobid {
            jobret = api.jobs[i]
            found = 1
            break
        }
    }
    api.mutex.Unlock()
    if found == 0 {
        return jobret, ApiError{}
    }
    return jobret, nil
}

func (api *Api) SetUserCancel( jobid int64 ) {
    api.mutex.Lock()
    for i,job := range api.jobs {
        if job.JobID == jobid {
            api.jobs[i].UserCancel = true
            break
        }
    }
    api.mutex.Unlock()
}

func (api *Api) UserCancel( jobid int64 ) bool {
    api.mutex.Lock()
    for i,job := range api.jobs {
        if job.JobID == jobid {
            if api.jobs[i].UserCancel == true {
                api.mutex.Unlock()
                return true;
            }
            break
        }
    }
    api.mutex.Unlock()
    return false
}

func (api *Api) SetPid( jobid int64, pid int64 ) {
    api.mutex.Lock()
    for i,job := range api.jobs {
        if job.JobID == jobid {
            api.jobs[i].Pid = pid
            break
        }
    }
    api.mutex.Unlock()
}

func (api *Api) RemoveJob( jobid int64 ) {
    api.mutex.Lock()
    i := -1
    for j,job := range api.jobs {
        if job.JobID == jobid {
            i = j
            break
        }
    }
    if i != -1 {
        api.jobs = append(api.jobs[:i], api.jobs[i+1:]...)
    }
    api.mutex.Unlock()
}

func (api *Api) Guid() string {
    api.mutex.Lock()
    defer api.mutex.Unlock()
    return api.guid
}

func (api *Api) UpdateGuid( guid string ) {
    api.mutex.Lock()
    api.guid = guid
    api.mutex.Unlock()
}

func (api *Api) Jobs( ) []JobIn {
    api.mutex.Lock()
    defer api.mutex.Unlock()
    return api.jobs
}

func NewApi() Api {
    api := Api{}
    api.jobs = make([]JobIn, 0)
    api.mutex = &sync.Mutex{}
    api.loginmutex = &sync.Mutex{}

    // bigger values allow for more concurrency.
    // increase to avoid 'use of closed network connection' errors
    if config.TransportTimeout < 2 {
        config.TransportTimeout = 2
    }

    return api
}

