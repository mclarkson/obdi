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
  "log"
  "log/syslog"
  "net"
  "net/http"
  "net/rpc"
  "crypto/tls"
  "encoding/json"
  "fmt"
  "time"
  "bytes"
  "strings"
  "io/ioutil"
  "os"
  "strconv"
)

const (
    SUCCESS = 0
    ERROR = 1
)

type ApiError struct {
  details string
}

func (e ApiError) Error() string {
  return fmt.Sprintf("%s", e.details)
}

// For sending a job to the Manager
type Job struct {
  Id            int64
  ScriptId      int64
  Args          string // E.g. `-a -f "bob 1" name`
  EnvVars       string // E.g. `A:1 B:"Hi there" C:3`
  Status        int64
  StatusReason  string
  StatusPercent int64
  CreatedAt     time.Time
  UpdatedAt     time.Time
  DeletedAt     time.Time
  UserLogin     string
  Errors        int64
  EnvId         int64 // For WorkerUrl and WorkerKey
  Type          int64 // 1 - user job, 2 - system job
}

// For retrieving details from the Manager
type Script struct {
    Id        int64
    Name      string
    Desc      string
    Source    []byte
    Type      string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time
}

// Args are send over RPC from the Manager
type Args struct {
  PathParams    map[string]string
  QueryString   map[string][]string
  PostData      []byte
  QueryType     string
}

// Log to syslog
func logit(msg string) {
    log.Println(msg)
    l, err := syslog.New(syslog.LOG_ERR, "obdi")
    defer l.Close()
    if err != nil {
        log.Fatal("error writing syslog!")
    }

    l.Err(msg)
}

/*
 * Send HTTP GET request
 */
func GET(url, endpoint string) (r *http.Response, e error) {

  // accept bad certs
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  client.Timeout = 8 * 1e9

  //fmt.Printf("\n%s/api/%s\n",url,endpoint)
  for strings.HasSuffix(url, "/") {
    url = strings.TrimSuffix(url, "/")
  }
  //fmt.Printf( "%s\n", url+"/"+endpoint )
  resp, err := client.Get(url+"/"+endpoint)
  if err != nil {
    txt := fmt.Sprintf("Could not send REST request ('%s').", err.Error())
    return resp, ApiError{txt}
  }

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
  }

  return resp, nil
}

/*
 * Send HTTP POST request
 */
func POST(jsondata []byte, url, endpoint string) (r *http.Response, e error) {

  buf := bytes.NewBuffer(jsondata)

  // accept bad certs
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  client.Timeout = 8 * 1e9

  //fmt.Printf("\n%s/api/%s\n",url,endpoint)
  for strings.HasSuffix(url, "/") {
    url = strings.TrimSuffix(url, "/")
  }
  //fmt.Printf( "%s\n", url+"/"+endpoint )
  resp, err := client.Post(url+"/"+endpoint, "application/json", buf)
  if err != nil {
    txt := fmt.Sprintf("Could not send REST request ('%s').", err.Error())
    return resp, ApiError{txt}
  }

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
  }

  return resp, nil
}

type Plugin struct{}

// The reply will be sent and output by the master
type Reply struct {
// Add more if required
JobId           int64
// Must have the following
PluginReturn    int64        // 0 - success, 1 - error
PluginError     string
}

// --------------------------------------------------------------------------
func ReturnError(text string, response *[]byte) {
// --------------------------------------------------------------------------
    errtext := Reply{ 0, ERROR, text }
    logit( text )
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
}

// --------------------------------------------------------------------------
func (t *Plugin) GetRequest(args *Args, response *[]byte) error {
// --------------------------------------------------------------------------

  // Check for required query string entries

  if len(args.QueryString["env_id"]) == 0 {
    txt := "'env_id' must be set"
    errtext := Reply{ 0,ERROR,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }

  env_id, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

  if len(args.QueryString["salt_id"]) == 0 {
    txt := "'salt_id' must be set"
    errtext := Reply{ 0,ERROR,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }

  // Get the ScriptId from the scripts table for:
  scriptName := "salt-grains-cache.sh"
  scripts := []Script{}
  resp, err := GET("https://127.0.0.1/api/" +
    args.PathParams["login"] + "/" + args.PathParams["GUID"], "scripts" +
    "?nosource=1&name=" + scriptName )
  if b, err := ioutil.ReadAll(resp.Body); err != nil {
    txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
    errtext := Reply{ 0,ERROR,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  } else {
    json.Unmarshal(b,&scripts)
  }
  // If scripts is empty then we don't have permission to see it
  // or the script does not exist (well, scripts don't have permissions
  // but lets say the same thing anyway)
  if len(scripts) == 0 {
    txt := "The requested script, '" + scriptName + "', does not exist" +
      " or the permissions to access it are insufficient."
    errtext := Reply{ 0,ERROR,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }

  // Set up some fields for the Job struct we'll send to the master
  job := Job{
    ScriptId:         scripts[0].Id,
    EnvId:            env_id,
    Args:             args.QueryString["salt_id"][0],

    // Type 1 - User Job - Output is
    //     sent back as it's created
    // Type 2 - System Job - All output
    //     is saved in one single line.
    //     Good for json etc.
    Type:             2,
  }

  // Send the job POST request to the master
  jsonjob, err := json.Marshal(job)
  resp, err = POST(jsonjob, "https://127.0.0.1/api/" +
    args.PathParams["login"] + "/" + args.PathParams["GUID"], "jobs")
  if err != nil {
    txt := "Could not send job to worker. ('" + err.Error() + "')"
    errtext := Reply{ 0,ERROR,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }
  defer resp.Body.Close()
  // Read the worker's response from the master
  if b, err := ioutil.ReadAll(resp.Body); err != nil {
    txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
    jsondata, _ := json.Marshal( txt )
    *response = jsondata
    return nil
  } else {
    json.Unmarshal(b,&job)
  }

  // Send the Job ID as the RPC reply (back to the master)

  id := job.Id
  reply := Reply{ id,SUCCESS,"" }
  jsondata, err := json.Marshal(reply)

  if err != nil {
    errtext := Reply{ 0,ERROR,"Marshal error: "+err.Error() }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }

  *response = jsondata

  return nil
}

// --------------------------------------------------------------------------
func (t *Plugin) PostRequest(args *Args, response *[]byte) error {
// --------------------------------------------------------------------------
    ReturnError( "Internal error: Unimplemented HTTP POST", response )
    return nil
}

/*
 * HandleRequest
 * All plugins must have this.
 */
func (t *Plugin) HandleRequest(args *Args, response *[]byte) error {

  if len(args.QueryType) > 0 {
    switch args.QueryType {
      case "GET": t.GetRequest(args, response)
                  return nil
      case "POST": t.PostRequest(args, response)
                   return nil
    }
    ReturnError( "Internal error: Invalid HTTP request type for this plugin " +
      args.QueryType, response )
    return nil
  } else {
    ReturnError( "Internal error: HTTP request type was not set", response )
    return nil
  }
}

func main() {

  //logit("Plugin starting")

  plugin := new(Plugin)
  rpc.Register(plugin)

  listener, err := net.Listen("tcp", ":" + os.Args[1])
  if err != nil {
      txt := fmt.Sprintf( "Listen error. ", err )
      logit( txt )
  }

  //logit("Plugin listening on port " + os.Args[1])

  //for {
    if conn, err := listener.Accept(); err != nil {
      txt := fmt.Sprintf( "Accept error. ", err )
      logit( txt )
    } else {
      //logit("New connection established")
      rpc.ServeConn(conn)
    }
  //}
}
