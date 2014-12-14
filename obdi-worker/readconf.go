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
	TransportTimeout int64  `toml:"transport_timeout"`
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
