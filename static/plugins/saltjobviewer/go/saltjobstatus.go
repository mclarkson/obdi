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
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/syslog"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

// ***************************************************************************
// SQLITE3 PRIVATE DB
// ***************************************************************************

type JobStatus struct {
	Id        int64
	JobId     string

  // Status
  // 0 - Unknown
  // 1 - Success
  // 2 - Fail
  // 3 - No output
	Status    int64
	CreatedAt time.Time
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

	gormInst.db, err = gorm.Open("sqlite3", dbname+"jobstatus.db")
	if err != nil {
		return ApiError{"Open " + dbname + " failed. " + err.Error()}
	}

	if err := gormInst.db.AutoMigrate(JobStatus{}).Error; err != nil {
		txt := fmt.Sprintf("AutoMigrate JobStatus table failed: %s", err)
		return ApiError{txt}
	}

	// Unique index is also a constraint, so are forced to be unique
	gormInst.db.Model(JobStatus{}).AddIndex("idx_jobstatus_job_id", "job_id")

	// TODO: disable this
	gormInst.db.LogMode(true)

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

// PortLock is a locker which locks by binding to a port on the loopback IPv4 interface
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

// // ***************************************************************************
// // FILE LOCKING
// // ***************************************************************************
//
// // FLock is a file-based lock
// type FLock struct {
//   fh *os.File
// }
//
// // NewFLock creates new Flock-based lock (unlocked first)
// func NewFLock(path string) (FLock, error) {
//   //os.Create(path)
//   fh, err := os.Open(path)
//   if err != nil {
//     return FLock{}, err
//   }
//   return FLock{fh: fh}, nil
// }
//
// // Lock acquires the lock, blocking
// func (lock FLock) Lock() error {
//   return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX)
// }
//
// // TryLock acquires the lock, non-blocking
// func (lock FLock) TryLock() (bool, error) {
//   err := syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
//   switch err {
//   case nil:
//     return true, nil
//   case syscall.EWOULDBLOCK:
//     return false, nil
//   }
//   return false, err
// }
//
// // Unlock releases the lock
// func (lock FLock) Unlock() error {
//   lock.fh.Close()
//   return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_UN)
// }

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
	JobId       string
	Status      int64
  KeepJobs    int64
}

type Plugin struct{}

// The reply will be sent and output by the master
type Reply struct {
	// Add more if required
	JobStatus string
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

	var err error

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

	db := gormInst.DB() // shortcut

	jobs := []JobStatus{}
  // No results is not an error
  Lock()
  dberr := db.Order("id desc").Find(&jobs)
  Unlock()
  if dberr.Error != nil {
    if !dberr.RecordNotFound() {
      ReturnError(dberr.Error.Error(), response)
      return nil
    }
  }

	// Create a slice of maps from users struct
	// to selectively copy database fields for display

	u := make([]map[string]interface{}, len(jobs))
	for i := range jobs {
		u[i] = make(map[string]interface{})
		u[i]["Id"] = jobs[i].Id
		u[i]["JobId"] = jobs[i].JobId
		u[i]["Status"] = jobs[i].Status
		u[i]["CreatedAt"] = jobs[i].CreatedAt
	}

	// Send the added record back

  jobstatus_tmp, err := json.Marshal(u)

  if err != nil {
    ReturnError("Marshal error: "+err.Error(), response)
    return nil
  }

  jobstatus_string := string(jobstatus_tmp)

  // Put the record in the Reply

	reply := Reply{jobstatus_string, SUCCESS, ""}
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

	//ReturnError( "Internal error: Unimplemented HTTP POST with data " +
	//  fmt.Sprintf(": %s",args.PostData), response )
	//return nil

	var err error

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

	var postedData PostedData

	if err := json.Unmarshal(args.PostData, &postedData); err != nil {
		txt := fmt.Sprintf("Error decoding JSON ('%s')"+".", err.Error())
		ReturnError("Error decoding the POST data ("+
			fmt.Sprintf("%s", args.PostData)+"). "+txt, response)
		return nil
	}

  // Salt JobIds are always 20 characters long
  if len(postedData.JobId) != 20 {
		ReturnError("Invalid JobId.", response)
		return nil
	}

	db := gormInst.DB() // shortcut

	// Add the JobId and Status

  jobstatus := JobStatus{
    JobId:     postedData.JobId,
    Status:    postedData.Status,
  }
  Lock()
  if err := db.Create(&jobstatus); err.Error != nil {
    Unlock()
    ReturnError(err.Error.Error(), response)
    return nil
  }
  Unlock()

  // Delete any jobs older than keep_jobs/24 days

  keep_jobs := postedData.KeepJobs;
  if keep_jobs > 0 {
    // delete from job_status where created_at < date('now','-8 days');
    days := (keep_jobs/24)+1
    txt := fmt.Sprintf("created_at < date('now','-%d days')", days)
    Lock()
    if err := db.Where(txt).Delete(JobStatus{}); err.Error != nil {
      Unlock()
      ReturnError(err.Error.Error(), response)
      return nil
    }
    Unlock()
  }

	// Send the added record back

  jobstatus_tmp, err := json.Marshal(jobstatus)

  if err != nil {
    ReturnError("Marshal error: "+err.Error(), response)
    return nil
  }

  jobstatus_string := string(jobstatus_tmp)

  // Put the record in the Reply

	reply := Reply{jobstatus_string, SUCCESS, ""}
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
	config.Port = 49994
	config.Portlock = NewPortLock(config.Port)

	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		txt := fmt.Sprintf("Listen error. ", err)
		logit(txt)
	}

	plugin := new(Plugin)
	rpc.Register(plugin)

	//for {
	if conn, err := listener.Accept(); err != nil {
		txt := fmt.Sprintf("Accept error. ", err)
		logit(txt)
	} else {
		rpc.ServeConn(conn)
	}
	//}
}

// vim:ts=2:sw=2:et
