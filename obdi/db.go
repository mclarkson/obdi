// Copyright 2016 Mark Clarkson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/mclarkson/obdi/external/jinzhu/gorm"
	_ "github.com/mclarkson/obdi/external/mattn/go-sqlite3"
	"log"
	"time"
)

// Maps to the users table
type User struct {
	Id         int64
	Login      string `sql:"not null"`
	Forename   string
	Surname    string
	Passhash   string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Email      string
	MultiLogin bool

	//Session     Session
	//SessionId   sql.NullInt64
	//Uemap     []Uemap
}

// Maps to the sessions table
type Session struct {
	Id        int64
	Guid      string
	UserId    int64 `sql:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Application log
type Activity struct {
	Id         int64
	Session_id int64
	Message    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

// Data Centre environment, e.g. Peer1 - Dev, or Newman St. - Prod
type Env struct {
	Id       int64
	DispName string // Display name
	SysName  string // System name (Salt environment name)
	/*Dc          Dc*/ // only for creating Env and substruct
	DcId               int64
	//WorkerIp    string      // Hostname or IP address of worker
	//WorkerPort  string      // Port the worker listens on
	WorkerUrl string // Worker URL Prefix
	WorkerKey string // Key (password) for worker
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Dc struct {
	Id        int64
	DispName  string // Display name
	SysName   string // System name (Salt name)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type DcCap struct {
	Id   int64
	Code string
	Desc string
}

type DcCapMap struct {
	Id      int64
	DcId    int64
	DcCapId int64
}

type EnvCap struct {
	Id              int64
	Code            string
	Desc            string
	IsWorkerDef     bool // Can it have a worker definition
	IsJsonObjectDef bool // Can it have a json object definition
}

type Worker struct {
	Id        int64
	EnvId     int64
	EnvCapId  int64
	WorkerUrl string // Worker URL Prefix
	WorkerKey string // Key (password) for worker
}

type JsonObject struct {
	Id       int64
	EnvId    int64
	EnvCapId int64
	Json     string
}

type EnvCapMap struct {
	Id       int64
	EnvId    int64
	EnvCapId int64
}

// Permissions
type Perm struct {
	Id        int64
	UserId    int64
	EnvId     int64
	Writeable bool
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

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
	EnvId         int64  // For WorkerUrl and WorkerKey, overridden by:
	EnvCapDesc    string // E.g. "EZBACKUP_WORKER"
	Type          int64  // 1 - user job, 2 - system job
	WorkerUrl     string // Worker URL Prefix
	WorkerKey     string // Key (password) for worker
}

type OutputLine struct {
	Id     int64
	Serial int64
	JobId  int64
	Text   string
	//Type            int64       // 0 - output, 1 - error output
}

type Plugin struct {
	Id        int64
	Name      string
	Desc      string
	Parent    string // Parents Name
	HasView   int64  // 1 - has a view, 2 - has not
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type File struct {
	Id       int64
	Name     string
	Desc     string
	Url      string
	Type     int64 // 1 - JS, 2 - HTML, 3 - GO, 4 - CSS
	PluginId int64
}

type Repo struct {
	Id        int64
	Url       string
	CreatedAt time.Time
}

type Database struct {
	dB gorm.DB
}

func (db *Database) InitDB() {
	var err error
	dbname := config.DBPath()

	db.dB, err = gorm.Open("sqlite3", dbname)
	if err != nil {
		logit("Could not open sqlite3 database, '" + dbname + "'")
		log.Fatal(fmt.Sprintf(
			"Got error when connect database, the error is '%v'", err))
	}

	if err := db.dB.AutoMigrate(User{}).Error; err != nil {
		txt := "AutoMigrate Users table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Session{}).Error; err != nil {
		txt := "AutoMigrate Sessions table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Activity{}).Error; err != nil {
		txt := "AutoMigrate Activity table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Env{}).Error; err != nil {
		txt := "AutoMigrate Env table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Dc{}).Error; err != nil {
		txt := "AutoMigrate Dc table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Perm{}).Error; err != nil {
		txt := "AutoMigrate Dc table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(DcCap{}).Error; err != nil {
		txt := "AutoMigrate DcCap table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(DcCapMap{}).Error; err != nil {
		txt := "AutoMigrate DcCapMap table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(EnvCap{}).Error; err != nil {
		txt := "AutoMigrate EnvCap table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(EnvCapMap{}).Error; err != nil {
		txt := "AutoMigrate EnvCapMap table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Perm{}).Error; err != nil {
		txt := "AutoMigrate Dc table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Job{}).Error; err != nil {
		txt := "AutoMigrate Job table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(OutputLine{}).Error; err != nil {
		txt := "AutoMigrate OutputLine table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Script{}).Error; err != nil {
		txt := "AutoMigrate Script table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Plugin{}).Error; err != nil {
		txt := "AutoMigrate Plugin table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(File{}).Error; err != nil {
		txt := "AutoMigrate File table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Repo{}).Error; err != nil {
		txt := "AutoMigrate Repos table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(Worker{}).Error; err != nil {
		txt := "AutoMigrate Workers table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}
	if err := db.dB.AutoMigrate(JsonObject{}).Error; err != nil {
		txt := "AutoMigrate JsonObjects table failed"
		log.Fatal(fmt.Sprintf("%s: %s", txt, err))
	}

	// Unique index is also a constraint. So these are forced to be unique
	db.dB.Model(User{}).AddUniqueIndex("idx_login", "login")
	db.dB.Model(Plugin{}).AddIndex("idx_name", "name")
	db.dB.Model(Session{}).AddIndex("idx_user_id", "user_id")
	db.dB.Model(Activity{}).AddIndex("idx_session_id", "session_id")
	db.dB.Model(Script{}).AddIndex("idx_script_name", "name")
	db.dB.Model(Worker{}).AddIndex("idx_worker_env_id", "env_id")
	db.dB.Model(Worker{}).AddIndex("idx_worker_env_cap_id", "env_cap_id")
	db.dB.Model(JsonObject{}).AddIndex("idx_jsonobject_env_id", "env_id")
	db.dB.Model(JsonObject{}).AddIndex("idx_jsonobject_env_cap_id", "env_cap_id")
	// TODO: OutputLines table should be in a separate DB file if
	// TODO: performance drops.
	db.dB.Model(OutputLine{}).AddIndex("idx_id_serial", "job_id", "serial")

	logit("Sqlite3 database " + dbname + " opened")
}

func (db *Database) CreateAdminAccount() {

	user := User{}
	user.Login = "admin"
	user.Enabled = true

	if db.dB.Where("login = ?", "admin").
		First(&user).RecordNotFound() {

		c := &Crypt{}
		c.Pass = []byte("admin")
		c.Crypt()

		db.dB.Where(User{Login: "admin"}).
			Attrs(User{
				Passhash: string(c.Hash),
				Forename: "Admin",
				Surname:  "User",
				Email:    "admin@invalid",
			}).FirstOrCreate(&user)

		logit("Admin user created")
	}
}

func (db *Database) DB() *gorm.DB {
	return &db.dB
}

func NewDB() *Database {
	db := &Database{}
	db.InitDB()
	return db
}
