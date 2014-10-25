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
  "net/rpc"
  "encoding/json"
  "fmt"
  "os"
  "sync"
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

const DBFILE = "statedesc.db"

// ***************************************************************************
// SQLITE3 PRIVATE DB
// ***************************************************************************

type StateDesc struct {
  Id              int64
  FormulaName     string
  StateFileName   string
  Desc            string
}

// --

var config *Config

type Config struct {
    Dbname    string
}

// --------------------------------------------------------------------------
func (c *Config) DBPath() string {
// --------------------------------------------------------------------------
    return c.Dbname
}

// --------------------------------------------------------------------------
func (c *Config) SetDBPath( path string ) {
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

  gormInst.db, err = gorm.Open("sqlite3", dbname + DBFILE)
  if err != nil {
    return ApiError{"Open " + dbname + " failed. " + err.Error()}
  }

  if err := gormInst.db.AutoMigrate(StateDesc{}).Error; err != nil {
    txt := fmt.Sprintf("AutoMigrate StateDesc table failed: %s", err)
    return ApiError{ txt }
  }

  // Unique index is also a constraint, so are forced to be unique
  //gormInst.db.Model(Enc{}).AddIndex("idx_enc_salt_id", "salt_id")

  return nil
}

// --------------------------------------------------------------------------
func (gormInst *GormDB) DB() *gorm.DB {
// --------------------------------------------------------------------------
  return &gormInst.db
}

// --------------------------------------------------------------------------
func NewDB() (*GormDB,error) {
// --------------------------------------------------------------------------
  gormInst := &GormDB{}
  if err := gormInst.InitDB(); err != nil {
    return gormInst, err
  }
  return gormInst,nil
}

// ***************************************************************************
// ERRORS
// ***************************************************************************

const (
    SUCCESS = 0
    ERROR = 1
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
// GO RPC PLUGIN
// ***************************************************************************

// Args are send over RPC from the Manager
type Args struct {
  PathParams    map[string]string
  QueryString   map[string][]string
  PostData      []byte
  QueryType     string
}

type Plugin struct{}

// The reply will be sent and output by the master
type Reply struct {
  // Add more if required
  EncData         string
  // Must have the following
  PluginReturn    int64        // 0 - success, 1 - error
  PluginError     string
}

// Global mutex
var dbmutex = &sync.Mutex{}

// --------------------------------------------------------------------------
func ReturnError(text string, response *[]byte) {
// --------------------------------------------------------------------------
    errtext := Reply{ "", ERROR, text }
    logit( text )
    jsondata, _ := json.Marshal( errtext )
    *response = jsondata
}

// --------------------------------------------------------------------------
func (t *Plugin) GetRequest(args *Args, response *[]byte) error {
// --------------------------------------------------------------------------

  // Check for required query string entries

  var err error

  if len(args.PathParams["PluginDatabasePath"]) == 0 {
    ReturnError( "Internal Error: 'PluginDatabasePath' must be set",response )
    return nil
  }

  config.SetDBPath( args.PathParams["PluginDatabasePath"] )

  // Open/Create database
  var gormInst *GormDB
  if gormInst,err = NewDB(); err!=nil {
    txt := "GormDB open error for '" + config.DBPath() + DBFILE + "'. " +
           err.Error()
    ReturnError( txt, response )
    return nil
  }

  db := gormInst.DB()          // Instead of using gormInst.DB() everywhere

  stateDescs := []StateDesc{}

  dbmutex.Lock()
  if err := db.Find(&stateDescs); err.Error != nil {
      if !err.RecordNotFound() {
        dbmutex.Unlock()
        ReturnError( err.Error.Error(), response )
        return nil
      }
  }
  dbmutex.Unlock()

  TempEncData, err := json.Marshal(stateDescs)
  EncData := string( TempEncData )

  if err != nil {
    ReturnError( "Marshal error: "+err.Error(), response )
    return nil
  }

  // Reply with the EncData (back to the master)

  reply := Reply{ EncData,SUCCESS,"" }
  jsondata, err := json.Marshal(reply)

  if err != nil {
    ReturnError( "Marshal error: "+err.Error(), response )
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

// --------------------------------------------------------------------------
func (t *Plugin) HandleRequest(args *Args, response *[]byte) error {
// --------------------------------------------------------------------------
// All plugins must have this.

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

// ***************************************************************************
// ENTRY POINT
// ***************************************************************************

// --------------------------------------------------------------------------
func main() {
// --------------------------------------------------------------------------

  // Sets the global config var
  NewConfig()

  plugin := new(Plugin)
  rpc.Register(plugin)

  listener, err := net.Listen("tcp", ":" + os.Args[1])
  if err != nil {
      txt := fmt.Sprintf( "Listen error. ", err )
      logit( txt )
  }

  //for {
    if conn, err := listener.Accept(); err != nil {
      txt := fmt.Sprintf( "Accept error. ", err )
      logit( txt )
    } else {
      rpc.ServeConn(conn)
    }
  //}
}
