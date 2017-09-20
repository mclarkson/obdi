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
	"os"
	//"fmt"
	"github.com/mclarkson/obdi/external/BurntSushi/toml"
)

var config Config

type Config struct {
	Dbname            string `toml:"database_path"`
	PluginDbPath      string `toml:"plugin_database_path"`
	ListenAddr        string `toml:"listen_address"`
	SessionTimeout    int    `toml:"session_timeout"`
	StaticContent     string `toml:"static_content"`
	SyslogEnabled     bool   `toml:"syslog_enabled"`
	SSLEnabled        bool   `toml:"ssl_enabled"`
	SSLCertFile       string `toml:"ssl_cert"`
	SSLKeyFile        string `toml:"ssl_key"`
	GoPluginDir       string `toml:"go_plugin_dir"`
	GoPluginSource    string `toml:"go_plugin_source"`
	GoPluginPortStart int64  `toml:"go_plugin_port_start"`
	GoRoot            string `toml:"go_root"`
	CacheDir          string `toml:"cache_dir"`
	HttpProxy         string `toml:"http_proxy"`
	HttpsProxy        string `toml:"https_proxy"`
	TransportTimeout  int64  `toml:"transport_timeout"` // Not used
}

func init() {
	config = Config{}
	config.Read_config()
}

func (c Config) DBPath() string {
	return c.Dbname
}

func (c *Config) Read_config() {
	if _, err := toml.DecodeFile("/etc/obdi/obdi.conf",
		c); err != nil {
		logit(err.Error())
		//txt := "Unable to read configuration"
		//fmt.Printf( "%s: %s\n", txt, err )
		os.Exit(1)
	}
}
