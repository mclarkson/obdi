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
	"regexp"
	"strings"
	//"syscall" // for flock
	"strconv"
	"time"
	//"math"
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
	Classes     []string
	Dc          string
	Environment string
}

type Plugin struct{}

// The reply will be sent and output by the master
type Reply struct {
	// Add more if required
	EncData string
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

	// Check for required query string entries

	var err error

	//dc_sent := 0
	dc_name := "" // E.g. OFFICE
	//env_sent := 0
	env_name := ""    // E.g. test_1.122
	env_version := "" // E.g. test
	salt_id := ""

	if len(args.QueryString["dc"]) > 0 {
		dc_name = args.QueryString["dc"][0]
		//dc_sent = 1
	} else { // Added due to NOTE 1 above
		ReturnError("The dc must be set", response)
		return nil
	}

	if len(args.QueryString["version"]) > 0 {
		env_version = args.QueryString["version"][0]
		//env_sent = 1
	} else { // Added due to NOTE 1 above
		ReturnError("The version must be set", response)
		return nil
	}

	if len(args.QueryString["env"]) > 0 {
		env_name = args.QueryString["env"][0]
		//env_sent = 1
	} else { // Added due to NOTE 1 above
		ReturnError("The env must be set", response)
		return nil
	}

	if len(args.QueryString["salt_id"]) == 0 {
		ReturnError("'salt_id' must be set", response)
		return nil
	}

	salt_id = args.QueryString["salt_id"][0]

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

	// Get ENC formula's and state files from enc tables
	// Do we care who can get this information? I'm guessing 'no'.

	db := gormInst.DB() // shortcut
	encs := []Enc{}

	// Search the encs DB

	Lock()
	if err := db.Find(&encs, "salt_id = ? and dc = ? and env = ?",
		salt_id, dc_name, env_name); err.Error != nil {
		if !err.RecordNotFound() {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
	}
	Unlock()

	var encClasses []string
	var encEnvironment string

	customised := false

	if len(encs) == 0 {

		// ENC entry does not exist. Make one on the fly from regexes
		// and return it to the user. Also fully realise the enc entry
		// from the regex.

		// Get all regexes for this dc and env
		regexes := []Regex{}
		Lock()
		if err := db.Find(&regexes, "dc = ? and env = ?",
			dc_name, env_name); err.Error != nil {
			if !err.RecordNotFound() {
				Unlock()
				ReturnError(err.Error.Error(), response)
				return nil
			}
		}
		Unlock()

		// Apply regexes to salt id - see which regexes qualify

		matched := false

	i_loop:
		for i := range regexes {
			tryRegex, err := regexp.Compile(regexes[i].Regex)
			if err == nil {
				if tryRegex.MatchString(salt_id) {
					// Add classes to the EncData ??
					Lock()
					regexSlsMaps := []RegexSlsMap{}
					if err := db.Find(&regexSlsMaps,
						"regex_id = ?", regexes[i].Id); err.Error != nil {
						if !err.RecordNotFound() {
							Unlock()
							ReturnError(err.Error.Error(), response)
							return nil
						}
					}
					Unlock()
					for j := range regexSlsMaps {
						matched = true

						// Should end with '.sls' but strip just in case
						stripped := strings.TrimSuffix(regexSlsMaps[j].StateFile, ".sls")

						item := ""
						if len(stripped) > 0 {
							item = regexSlsMaps[j].Formula + "." + stripped
						} else {
							item = regexSlsMaps[j].Formula
						}
						// Check for existing entry
						for j := range encClasses {
							if encClasses[j] == item {
								continue i_loop
							}
						}

						encClasses = append(encClasses, item)

						// And also write to encs table

						enc := Enc{
							SaltId:    salt_id,
							Formula:   regexSlsMaps[j].Formula,
							StateFile: regexSlsMaps[j].StateFile,
							Dc:        dc_name,
							Env:       env_name,
						}
						Lock()
						if err := db.Create(&enc); err.Error != nil {
							Unlock()
							ReturnError(err.Error.Error(), response)
							return nil
						}
						Unlock()
					}
				}
			} else {
				logit(fmt.Sprintf("Regex error with '%s' (%s,%s,%s): %s.",
					regexes[i].Regex, dc_name, env_name, env_name+" "+env_version, err))
			}
		}

		// No regex either. Return empty handed

		if matched == false {

			logit("No classes found for " + salt_id +
				", and could not find a regex!")

		} else {
			encEnvironment = env_name // + "_" + env_version
		}

	} else {

		// ENC entry exists

		customised = true

	i_loop2:
		for i := range encs {

			stripped := strings.TrimSuffix(encs[i].StateFile, ".sls")

			item := ""
			if len(stripped) > 0 {
				item = encs[i].Formula + "." + stripped
			} else {
				item = encs[i].Formula
			}
			// Check for existing entry
			for j := range encClasses {
				if encClasses[j] == item {
					continue i_loop2
				}
			}
			encClasses = append(encClasses, item)
		}
		encEnvironment = env_name // + "_" + env_version
	}

	// Output as JSON or YAML

	var EncData string

	// Differences between data included in yaml and json output.
	//
	// yaml: Sends "environment_version", so salt can match the git tag/branch.
	//
	// json: Send environment name only.
	//       Includes 'customised' field. So GUI can tell where the data was
	//       sourced from - enc or regex.

	if len(args.QueryString["yaml"]) > 0 {

		// Output as YAML

		if len(encClasses) == 0 {
			EncData = "classes: null"
		} else {
			EncData = "classes:\n"
			for i := range encClasses {
				EncData += "  - " + encClasses[i] + "\n"
				// TODO: Salt enc does not support parameters
				//EncData += "  " + encClasses[i] + ":\n"
				//if( encClasses[i] == "vim.vimrc" ) {
				//  EncData += "    vim_new_value: 1\n"
				//}
			}
			// TODO: Salt enc does not support parameters
			//EncData += "parameters:\n"
			//EncData += "  vim_test_value: 2\n"
			EncData += "environment: " + encEnvironment + "_" + env_version
		}
	} else {

		// Output as JSON

		type JsonOut struct {
			Classes     []string
			Environment string
			Customised  bool
		}

		jsonout := JsonOut{encClasses, encEnvironment, customised}

		TempEncData, err := json.Marshal(jsonout)

		if err != nil {
			ReturnError("Marshal error: "+err.Error(), response)
			return nil
		}

		EncData = string(TempEncData)
	}

	// Reply with the EncData (back to the master)

	reply := Reply{EncData, SUCCESS, ""}
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

	if len(args.QueryString["salt_id"]) == 0 {
		ReturnError("'salt_id' must be set", response)
		return nil
	}

	salt_id := args.QueryString["salt_id"][0]

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

	// Remove all ENC Classes (before adding)

	db := gormInst.DB() // shortcut

	Lock()
	if err := db.Where("salt_id = ? and dc = ? and env = ?",
		salt_id, postedData.Dc, postedData.Environment).Delete(Enc{}); err.Error != nil {
		if !err.RecordNotFound() {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
	}
	Unlock()

	// Add the ENC classes

	for i := range postedData.Classes {
		classes := strings.Split(postedData.Classes[i], ".")
		formula := ""
		statefile := ""
		switch len(classes) {
		case 0:
			continue
		case 1:
			formula = classes[0]
		case 2:
			formula = classes[0]
			statefile = classes[1]
		}
		enc := Enc{
			SaltId:    salt_id,
			Formula:   formula,
			StateFile: statefile,
			Dc:        postedData.Dc,
			Env:       postedData.Environment,
		}
		Lock()
		if err := db.Create(&enc); err.Error != nil {
			Unlock()
			ReturnError(err.Error.Error(), response)
			return nil
		}
		Unlock()
	}

	reply := Reply{"", SUCCESS, ""}
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
	config.Port = 49993
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
