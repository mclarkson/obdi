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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	//"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

// ***************************************************************************
// SQLITE3 PRIVATE DB
// ***************************************************************************

type Enc struct {
	Id        int64
	SaltId    string // Name of the server
	Formula   string // Directory name
	StateFile string // Sls file name
	Dc        string // Data centre name
	Env       string // Environment name
}

type Regex struct {
	Id    int64
	Regex string // The regular expression
	Dc    string // Data centre name
	Env   string // Environment name
	Name  string // Short name for the regex, no spaces
	Desc  string // Description of the regex
}

type RegexSlsMap struct {
	Id        int64
	RegexId   int64  // Not null
	Formula   string // Not null
	StateFile string // Can be null
}

// For retrieving details from the Manager
type Env struct {
	Id       int64
	DispName string // Display name
	SysName  string // System name (Salt name)
	/*Dc          Dc*/ // only for creating Env and substruct
	DcId               int64
	DcSysName          string
	//WorkerIp    string      // Hostname or IP address of worker
	//WorkerPort  string      // Port the worker listens on
	WorkerUrl string // Worker URL Prefix
	WorkerKey string // Key (password) for worker
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// --

var config *Config

type Config struct {
	Dbname   string
	Portlock *PortLock
	Port     int
}

// --------------------------------------------------------------------------
func (c *Config) DBPath() string {
	// --------------------------------------------------------------------------
	return c.Dbname
}

// --------------------------------------------------------------------------
func (c *Config) SetDBPath(path string) {
	// --------------------------------------------------------------------------
	c.Dbname = path
}

// --------------------------------------------------------------------------
func NewConfig() {
	// --------------------------------------------------------------------------
	config = &Config{}
}

// --

type GormDB struct {
	db gorm.DB
}

// --------------------------------------------------------------------------
func (gormInst *GormDB) InitDB() error {
	// --------------------------------------------------------------------------
	var err error
	dbname := config.DBPath()

	gormInst.db, err = gorm.Open("sqlite3", dbname+"enc.db")
	if err != nil {
		return ApiError{"Open " + dbname + " failed. " + err.Error()}
	}

	if err := gormInst.db.AutoMigrate(Enc{}).Error; err != nil {
		txt := fmt.Sprintf("AutoMigrate Enc table failed: %s", err)
		return ApiError{txt}
	}
	if err := gormInst.db.AutoMigrate(Regex{}).Error; err != nil {
		txt := fmt.Sprintf("AutoMigrate Regex table failed: %s", err)
		return ApiError{txt}
	}
	if err := gormInst.db.AutoMigrate(RegexSlsMap{}).Error; err != nil {
		txt := fmt.Sprintf("AutoMigrate RegexSlsMap table failed: %s", err)
		return ApiError{txt}
	}

	// Unique index is also a constraint, so are forced to be unique
	gormInst.db.Model(Enc{}).AddIndex("idx_enc_salt_id", "salt_id")

	return nil
}

// --------------------------------------------------------------------------
func (gormInst *GormDB) DB() *gorm.DB {
	// --------------------------------------------------------------------------
	return &gormInst.db
}

// --------------------------------------------------------------------------
func NewDB() (*GormDB, error) {
	// --------------------------------------------------------------------------
	gormInst := &GormDB{}
	if err := gormInst.InitDB(); err != nil {
		return gormInst, err
	}
	return gormInst, nil
}

// ***************************************************************************
// ERRORS
// ***************************************************************************

const (
	SUCCESS = 0
	ERROR   = 1
)

type ApiError struct {
	details string
}

// --------------------------------------------------------------------------
func (e ApiError) Error() string {
	// --------------------------------------------------------------------------
	return fmt.Sprintf("%s", e.details)
}

// ***************************************************************************
// LOGGING
// ***************************************************************************

// --------------------------------------------------------------------------
func logit(msg string) {
	// --------------------------------------------------------------------------
	// Log to syslog
	log.Println(msg)
	l, err := syslog.New(syslog.LOG_ERR, "obdi")
	defer l.Close()
	if err != nil {
		log.Fatal("error writing syslog!")
	}

	l.Err(msg)
}

// ***************************************************************************
// PORT LOCKING
// ***************************************************************************

// PortLock is a locker which locks by binding to a port on the loopback IPv4
// interface
type PortLock struct {
	hostport string
	ln       net.Listener
}

// --------------------------------------------------------------------------
func NewPortLock(port int) *PortLock {
	// --------------------------------------------------------------------------
	// NewFLock creates new Flock-based lock (unlocked first)
	return &PortLock{hostport: net.JoinHostPort("127.0.0.1", strconv.Itoa(port))}
}

// --------------------------------------------------------------------------
func (p *PortLock) Lock() {
	// --------------------------------------------------------------------------
	// Lock acquires the lock, blocking
	t := 50 * time.Millisecond
	for {
		if l, err := net.Listen("tcp", p.hostport); err == nil {
			p.ln = l // thanks to zhangpy
			return
		}
		//log.Printf("spinning lock on %s (%s)", p.hostport, err)
		time.Sleep(t)
		//t = time.Duration(
		//  math.Min( float64(time.Duration(float32(t) * 1.5)), 2000 ))
	}
}

// --------------------------------------------------------------------------
func (p *PortLock) Unlock() {
	// --------------------------------------------------------------------------
	// Unlock releases the lock
	if p.ln != nil {
		p.ln.Close()
	}
}

// ***************************************************************************
// REST SUPPORT FUNCS
// ***************************************************************************

// --------------------------------------------------------------------------
func GET(url, endpoint string) (r *http.Response, e error) {
	// --------------------------------------------------------------------------
	// Send HTTP GET request

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
	resp, err := client.Get(url + "/" + endpoint)
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

// --------------------------------------------------------------------------
func POST(jsondata []byte, url, endpoint string) (r *http.Response, e error) {
	// --------------------------------------------------------------------------
	// Send HTTP POST request

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

// ***************************************************************************
// GO RPC PLUGIN
// ***************************************************************************

// Args are send over RPC from the Manager
type Args struct {
	PathParams  map[string]string
	QueryString map[string][]string
	PostData    []byte
	QueryType   string
}

type PostedData struct {
	// Dc and Env are retrieved from the env_id
	//Dc            string
	//Env           string
	Desc  string
	Id    int64
	Name  string
	Regex string
}

type Plugin struct{}

// The reply will be sent and output by the master
type Reply struct {
	// Add more if required
	JsonData string
	// Must have the following
	PluginReturn int64 // 0 - success, 1 - error
	PluginError  string
}

// --------------------------------------------------------------------------
func Unlock() {
	// --------------------------------------------------------------------------
	config.Portlock.Unlock()
}

// --------------------------------------------------------------------------
func Lock() {
	// --------------------------------------------------------------------------
	config.Portlock.Lock()
}

// --------------------------------------------------------------------------
func ReturnError(text string, response *[]byte) {
	// --------------------------------------------------------------------------
	errtext := Reply{"", ERROR, text}
	logit(text)
	jsondata, _ := json.Marshal(errtext)
	*response = jsondata
}

// --------------------------------------------------------------------------
func (t *Plugin) GetRequest(args *Args, response *[]byte) error {
	// --------------------------------------------------------------------------
	// Return list of all regexes for an environment

	// Check for required query string entries

	var err error

	if len(args.QueryString["env_id"]) == 0 {
		ReturnError("'env_id' must be set", response)
		return nil
	}

	env_id := args.QueryString["env_id"][0]
	//env_id_str, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

	// Get the Dc (DcSysName) and Env (SysName) for this env_id using REST.
	// The Data Centre name and Environment name are stored in:
	//   envs[0].DcSysName and envs[0].SysName
	// GET queries always return an array of items, even for 1 item.
	envs := []Env{}
	resp, err := GET("https://127.0.0.1/api/"+
		args.PathParams["login"]+"/"+args.PathParams["GUID"], "envs"+
		"?env_id="+env_id)
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	} else {
		json.Unmarshal(b, &envs)
	}
	// If envs is empty then we don't have permission to see it
	// or the env does not exist so bug out.
	if len(envs) == 0 {
		txt := "The requested environment id does not exist" +
			" or the permissions to access it are insufficient."
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	}

	dc := envs[0].DcSysName
	env := envs[0].SysName

	// PluginDatabasePath is required to open our private db
	if len(args.PathParams["PluginDatabasePath"]) == 0 {
		ReturnError("Internal Error: 'PluginDatabasePath' must be set", response)
		return nil
	}

	config.SetDBPath(args.PathParams["PluginDatabasePath"])

	// Open/Create database
	var gormInst *GormDB
	if gormInst, err = NewDB(); err != nil {
		txt := "GormDB open error for '" + config.DBPath() + "enc.db'. " +
			err.Error()
		ReturnError(txt, response)
		return nil
	}

	// Get Regex formula's and state files from enc tables
	// Do we care who can get this information? I'm guessing 'no'.

	db := gormInst.DB() // shortcut

	// Search the regexes table

	regexes := []Regex{}
	Lock()
	if err := db.Find(&regexes, "dc = ? and env = ?", dc, env); err.Error != nil {
		if !err.RecordNotFound() {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
	}
	Unlock()

	//   // Output as JSON

	u := make([]map[string]interface{}, len(regexes))
	for i := range regexes {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = regexes[i].Id
		u[i]["Regex"] = regexes[i].Regex
		u[i]["Dc"] = regexes[i].Dc
		u[i]["Env"] = regexes[i].Env
		u[i]["Name"] = regexes[i].Name
		u[i]["Desc"] = regexes[i].Desc
	}

	//type JsonOut struct {
	//  Text     string
	//}

	TempJsonData, err := json.Marshal(u)
	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}
	reply := Reply{string(TempJsonData), SUCCESS, ""}
	jsondata, err := json.Marshal(reply)

	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}

	*response = jsondata

	return nil
}

// --------------------------------------------------------------------------
func (t *Plugin) PostRequest(args *Args, response *[]byte) error {
	// --------------------------------------------------------------------------

	var err error

	// Needed if the salt version has been changed
	if len(args.QueryString["env_id"]) == 0 {
		ReturnError("'env_id' must be set", response)
		return nil
	}

	env_id := args.QueryString["env_id"][0]
	//env_id_str, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

	// Get the Dc (DcSysName) and Env (SysName) for this env_id using REST.
	// The Data Centre name and Environment name are stored in:
	//   envs[0].DcSysName and envs[0].SysName
	// GET queries always return an array of items, even for 1 item.
	envs := []Env{}
	resp, err := GET("https://127.0.0.1/api/"+
		args.PathParams["login"]+"/"+args.PathParams["GUID"], "envs"+
		"?env_id="+env_id)
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	} else {
		json.Unmarshal(b, &envs)
	}
	// If envs is empty then we don't have permission to see it
	// or the env does not exist so bug out.
	if len(envs) == 0 {
		txt := "The requested environment id does not exist" +
			" or the permissions to access it are insufficient."
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	}

	dc := envs[0].DcSysName
	env := envs[0].SysName

	// PluginDatabasePath is required to open our private db
	if len(args.PathParams["PluginDatabasePath"]) == 0 {
		ReturnError("Internal Error: 'PluginDatabasePath' must be set", response)
		return nil
	}

	config.SetDBPath(args.PathParams["PluginDatabasePath"])

	// Open/Create database
	var gormInst *GormDB
	if gormInst, err = NewDB(); err != nil {
		txt := "GormDB open error for '" + config.DBPath() + "enc.db'. " +
			err.Error()
		ReturnError(txt, response)
		return nil
	}

	// Decode the post data into struct

	var postdata PostedData

	if err := json.Unmarshal(args.PostData, &postdata); err != nil {
		txt := fmt.Sprintf("Error decoding JSON ('%s')"+".", err.Error())
		ReturnError("Error decoding the POST data ("+
			fmt.Sprintf("%s", args.PostData)+"). "+txt, response)
		return nil
	}

	db := gormInst.DB() // shortcut

	// The following regex will be written to the db
	regex := Regex{
		0,
		postdata.Regex,
		dc,
		env,
		postdata.Name,
		postdata.Desc,
	}

	// Update the Regex entry

	Lock()
	if err := db.Save(&regex).Error; err != nil {
		Unlock()
		ReturnError("Update error: "+err.Error(), response)
		return nil
	}
	Unlock()

	// Output JSON

	//jsonout := JsonOut { "PutRequest" }
	TempJsonData, err := json.Marshal(regex)
	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}
	reply := Reply{string(TempJsonData), SUCCESS, ""}
	jsondata, err := json.Marshal(reply)

	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}

	*response = jsondata

	return nil
}

// --------------------------------------------------------------------------
func (t *Plugin) PutRequest(args *Args, response *[]byte) error {
	// --------------------------------------------------------------------------

	var err error

	// Needed if the salt version has been changed
	if len(args.QueryString["env_id"]) == 0 {
		ReturnError("'env_id' must be set", response)
		return nil
	}

	env_id := args.QueryString["env_id"][0]
	//env_id_str, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

	// Get the Dc (DcSysName) and Env (SysName) for this env_id using REST.
	// The Data Centre name and Environment name are stored in:
	//   envs[0].DcSysName and envs[0].SysName
	// GET queries always return an array of items, even for 1 item.
	envs := []Env{}
	resp, err := GET("https://127.0.0.1/api/"+
		args.PathParams["login"]+"/"+args.PathParams["GUID"], "envs"+
		"?env_id="+env_id)
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	} else {
		json.Unmarshal(b, &envs)
	}
	// If envs is empty then we don't have permission to see it
	// or the env does not exist so bug out.
	if len(envs) == 0 {
		txt := "The requested environment id does not exist" +
			" or the permissions to access it are insufficient."
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	}

	dc := envs[0].DcSysName
	env := envs[0].SysName

	// PluginDatabasePath is required to open our private db
	if len(args.PathParams["PluginDatabasePath"]) == 0 {
		ReturnError("Internal Error: 'PluginDatabasePath' must be set", response)
		return nil
	}

	config.SetDBPath(args.PathParams["PluginDatabasePath"])

	// Open/Create database
	var gormInst *GormDB
	if gormInst, err = NewDB(); err != nil {
		txt := "GormDB open error for '" + config.DBPath() + "enc.db'. " +
			err.Error()
		ReturnError(txt, response)
		return nil
	}

	// Decode the post data into struct

	var postdata PostedData

	if err := json.Unmarshal(args.PostData, &postdata); err != nil {
		txt := fmt.Sprintf("Error decoding JSON ('%s')"+".", err.Error())
		ReturnError("Error decoding the POST data ("+
			fmt.Sprintf("%s", args.PostData)+"). "+txt, response)
		return nil
	}

	db := gormInst.DB() // shortcut

	// Search the regexes table for the regex id

	regexes := []Regex{}
	Lock()
	if err := db.Find(&regexes, "id = ? and dc = ? and env = ?", postdata.Id,
		dc, env); err.Error != nil {
		if !err.RecordNotFound() {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
		id := strconv.FormatInt(postdata.Id, 10)
		ReturnError("Regex Id:"+id+" not found", response)
		return nil
	}
	Unlock()

	// The following regex will be written to the db
	regex := Regex{
		postdata.Id,
		postdata.Regex,
		dc,
		env,
		postdata.Name,
		postdata.Desc,
	}

	// Update the Regex entry

	Lock()
	if err := db.Save(&regex).Error; err != nil {
		Unlock()
		ReturnError("Update error: "+err.Error(), response)
		return nil
	}
	Unlock()

	// Output JSON

	//jsonout := JsonOut { "PutRequest" }
	TempJsonData, err := json.Marshal(regex)
	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}
	reply := Reply{string(TempJsonData), SUCCESS, ""}
	jsondata, err := json.Marshal(reply)

	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}

	*response = jsondata

	return nil
}

// --------------------------------------------------------------------------
func (t *Plugin) DeleteRequest(args *Args, response *[]byte) error {
	// --------------------------------------------------------------------------

	var err error

	// Needed if the salt version has been changed
	if len(args.QueryString["env_id"]) == 0 {
		ReturnError("'env_id' must be set", response)
		return nil
	}

	env_id := args.QueryString["env_id"][0] //string
	//env_id, _ := strconv.ParseInt( args.QueryString["env_id"][0],10,64 )

	// Get the Dc (DcSysName) and Env (SysName) for this env_id using REST.
	// The Data Centre name and Environment name are stored in:
	//   envs[0].DcSysName and envs[0].SysName
	// GET queries always return an array of items, even for 1 item.
	envs := []Env{}
	resp, err := GET("https://127.0.0.1/api/"+
		args.PathParams["login"]+"/"+args.PathParams["GUID"], "envs"+
		"?env_id="+env_id)
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	} else {
		json.Unmarshal(b, &envs)
	}
	// If envs is empty then we don't have permission to see it
	// or the env does not exist so bug out.
	if len(envs) == 0 {
		txt := "The requested environment id does not exist" +
			" or the permissions to access it are insufficient."
		errtext := Reply{"", ERROR, txt}
		jsondata, _ := json.Marshal(errtext)
		*response = jsondata
		return nil
	}

	dc := envs[0].DcSysName
	env := envs[0].SysName

	// PluginDatabasePath is required to open our private db
	if len(args.PathParams["PluginDatabasePath"]) == 0 {
		ReturnError("Internal Error: 'PluginDatabasePath' must be set", response)
		return nil
	}

	config.SetDBPath(args.PathParams["PluginDatabasePath"])

	// Open/Create database
	var gormInst *GormDB
	if gormInst, err = NewDB(); err != nil {
		txt := "GormDB open error for '" + config.DBPath() + "enc.db'. " +
			err.Error()
		ReturnError(txt, response)
		return nil
	}

	id_str := args.PathParams["id"]
	id_int, _ := strconv.ParseInt(args.PathParams["id"], 10, 64)

	db := gormInst.DB() // shortcut

	// Search the regexes table for the regex id

	regexes := []Regex{}
	Lock()
	if err := db.Find(&regexes, "id = ? and dc = ? and env = ?", id_str,
		dc, env); err.Error != nil {
		if !err.RecordNotFound() {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
		ReturnError("Regex Id:"+id_str+" not found", response)
		return nil
	}
	Unlock()

	// The following regex will be written to the db
	regex := Regex{
		id_int,
		"",
		dc,
		env,
		"",
		"",
	}

	// Update the Regex entry

	Lock()
	if err := db.Delete(&regex).Error; err != nil {
		Unlock()
		ReturnError("Update error: "+err.Error(), response)
		return nil
	}
	Unlock()

	// Output JSON

	//jsonout := JsonOut { "PutRequest" }
	TempJsonData, err := json.Marshal(regex)
	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}
	reply := Reply{string(TempJsonData), SUCCESS, ""}
	jsondata, err := json.Marshal(reply)

	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}

	*response = jsondata

	return nil
}

// --------------------------------------------------------------------------
func (t *Plugin) HandleRequest(args *Args, response *[]byte) error {
	// --------------------------------------------------------------------------
	// All plugins must have this.

	if len(args.QueryType) > 0 {
		switch args.QueryType {
		case "GET":
			t.GetRequest(args, response)
			return nil
		case "POST":
			t.PostRequest(args, response)
			return nil
		case "PUT":
			t.PutRequest(args, response)
			return nil
		case "DELETE":
			t.DeleteRequest(args, response)
			return nil
		}
		ReturnError("Internal error: Invalid HTTP request type for this plugin "+
			args.QueryType, response)
		return nil
	} else {
		ReturnError("Internal error: HTTP request type was not set", response)
		return nil
	}
}

// ***************************************************************************
// ENTRY POINT
// ***************************************************************************

// --------------------------------------------------------------------------
func main() {
	// --------------------------------------------------------------------------

	// Sets the global config var
	NewConfig()

	// Create a lock file to use for synchronisation
	config.Port = 49993
	config.Portlock = NewPortLock(config.Port)

	plugin := new(Plugin)
	rpc.Register(plugin)

	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		txt := fmt.Sprintf("Listen error. ", err)
		logit(txt)
	}

	//for {
	if conn, err := listener.Accept(); err != nil {
		txt := fmt.Sprintf("Accept error. ", err)
		logit(txt)
	} else {
		rpc.ServeConn(conn)
	}
	//}
}

// vim:ts=2:sw=2
