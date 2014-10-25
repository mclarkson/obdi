// server.go
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

type ApiError struct {
  details string
}

func (e ApiError) Error() string {
  return fmt.Sprintf("%s", e.details)
}

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

type Args struct {
  PathParams map[string]string
  QueryString map[string][]string
}

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

/*
 * HandleRequest
 * All plugins must have this.
 */
func (t *Plugin) HandleRequest(args *Args, response *[]byte) error {

  //fmt.Printf( "Processing request for %s\n", args.PathParams["login"] )

  //
  // Every plugin must return at least 'PluginReturn' and 'PluginError' fields
  //
  type Reply struct {
    JobId           int64
    PluginReturn    int64        // 0 - success, 1 - error
    PluginError     string
  }

  // Check for required query string entries

  if len(args.QueryString["script_id"]) == 0 {
    txt := "'script_id' must be set"
    errtext := Reply{ 0,1,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }
  script_id, _ := strconv.ParseInt( args.QueryString["script_id"][0],10,64 )

  if len(args.QueryString["env_id"]) == 0 {
    txt := "'env_id' must be set"
    errtext := Reply{ 0,1,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }
  env_id, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

  // Permissions check. Can 'login' access 'env_id'?

  //TODO: Permissions check on env_id

  // POST the job request for the worker to the master

  job := Job{
    ScriptId:         script_id,
    EnvId:            env_id,
    // TODO //Args:             dc and env from envtable,

    // Type 1 - User Job - Output is
    //     sent back as it's created
    // Type 2 - System Job - All output
    //     is saved in one single line.
    //     Good for json etc.
    Type:             2,
  }

  jsonjob, err := json.Marshal(job)
  resp, err := POST(jsonjob, "https://127.0.0.1/api/" +
    args.PathParams["login"] + "/" + args.PathParams["GUID"], "jobs")

  if err != nil {
    txt := "Could not send job to worker. ('" + err.Error() + "')"
    errtext := Reply{ 0,1,txt }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }
  defer resp.Body.Close()

  // Read the response for the worker from the master

  if b, err := ioutil.ReadAll(resp.Body); err != nil {
    txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
    jsondata, _ := json.Marshal( txt )
    *response = jsondata
    return nil
  } else {
    json.Unmarshal(b,&job)
  }

  // Send the Job ID as the RPC reply (back to the master)

  id := int64(job.Id)
  reply := Reply{ id,0,"" }
  jsondata, err := json.Marshal(reply)

  if err != nil {
    errtext := Reply{ 0,1,"Marshal error"+err.Error() }
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
    return nil
  }

  *response = jsondata

  return nil
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
