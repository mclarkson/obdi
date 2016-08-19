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
	"os"
	//"github.com/BurntSushi/toml"
	"github.com/mclarkson/obdi/external/BurntSushi/toml"
)

var config Config

type Config struct {
	Dbname           string `toml:"database_path"`
	ListenAddr       string `toml:"listen_address"`
	SessionTimeout   int    `toml:"session_timeout"`
	StaticContent    string `toml:"static_content"`
	SSLEnabled       bool   `toml:"ssl_enabled"`
	SSLCertFile      string `toml:"ssl_cert"`
	SSLKeyFile       string `toml:"ssl_key"`
	WorkerKey        string `toml:"key"`
	ScriptDir        string `toml:"script_dir"`
	User             string `toml:"man_user"`
	Password         string `toml:"man_password"`
	ManUrlPrefix     string `toml:"man_urlprefix"`
	SysScriptDir     string `toml:"system_scripts"`
	TransportTimeout int64  `toml:"transport_timeout"` // Not used
}

func init() {
	config = Config{}
	config.Read_config()
}

func (c *Config) Read_config() {
	if _, err := toml.DecodeFile("/etc/obdi-worker/obdi-worker.conf",
		c); err != nil {
		logit(err.Error())
		txt := "Unable to read configuration"
		fmt.Printf("%s: %s\n", txt, err)
		os.Exit(1)
	}
}
