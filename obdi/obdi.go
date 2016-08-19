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

// obdi manages deployment processes
package main

import (
	//"fmt"
	"crypto/tls"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}
var apimutex = &sync.Mutex{}

func main() {

	// Assign to global Transport var - accept self-signed certs
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	db := NewDB()
	api := NewApi(db)

	db.CreateAdminAccount()

	// TODO: disable this
	db.DB().LogMode(false)

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		//DisableJsonIndent: true,
		//EnableGzip:        true,
	}
	handler.SetRoutes(

		&rest.Route{"POST", "/login", api.DoLogin},

		&rest.Route{"POST", "/#login/:GUID/logout", api.Logout},

		// ADMIN FUNCTIONS

		// USERS

		&rest.Route{"GET", "/:login/:GUID/users", api.GetAllUsers},

		&rest.Route{"POST", "/:login/:GUID/users", api.AddUser},

		&rest.Route{"PUT", "/:login/:GUID/users/:id", api.UpdateUser},

		&rest.Route{"DELETE", "/:login/:GUID/users/:id", api.DeleteUser},

		// Data Centres

		&rest.Route{"GET", "/:login/:GUID/dcs", api.GetAllDcs},

		&rest.Route{"POST", "/:login/:GUID/dcs", api.AddDc},

		&rest.Route{"DELETE", "/:login/:GUID/dcs/:id", api.DeleteDc},

		&rest.Route{"PUT", "/:login/:GUID/dcs/:id", api.UpdateDc},

		// Environments

		&rest.Route{"GET", "/#login/:GUID/envs", api.GetAllEnvs},

		&rest.Route{"POST", "/:login/:GUID/envs", api.AddEnv},

		&rest.Route{"DELETE", "/:login/:GUID/envs/:id", api.DeleteEnv},

		&rest.Route{"PUT", "/:login/:GUID/envs/:id", api.UpdateEnv},

		// User Permissions

		&rest.Route{"GET", "/:login/:GUID/perms", api.GetAllPerms},

		&rest.Route{"POST", "/:login/:GUID/perms", api.AddPerm},

		&rest.Route{"DELETE", "/:login/:GUID/perms/:id", api.DeletePerm},

		&rest.Route{"PUT", "/:login/:GUID/perms/:id", api.UpdatePerm},

		// Data Centre Capabilities

		&rest.Route{"GET", "/:login/:GUID/dccaps", api.GetAllDcCaps},

		&rest.Route{"POST", "/:login/:GUID/dccaps", api.AddDcCap},

		&rest.Route{"DELETE", "/:login/:GUID/dccaps/:id", api.DeleteDcCap},

		&rest.Route{"PUT", "/:login/:GUID/dccaps/:id", api.UpdateDcCap},

		// Data Centre Capability Maps

		&rest.Route{"GET", "/:login/:GUID/dccapmaps", api.GetAllDcCapMaps},

		&rest.Route{"POST", "/:login/:GUID/dccapmaps", api.AddDcCapMap},

		&rest.Route{"DELETE", "/:login/:GUID/dccapmaps/:id", api.DeleteDcCapMap},

		&rest.Route{"PUT", "/:login/:GUID/dccapmaps/:id", api.UpdateDcCapMap},

		// Environment Capabilities

		&rest.Route{"GET", "/#login/:GUID/envcaps", api.GetAllEnvCaps},

		&rest.Route{"POST", "/:login/:GUID/envcaps", api.AddEnvCap},

		&rest.Route{"DELETE", "/:login/:GUID/envcaps/:id", api.DeleteEnvCap},

		&rest.Route{"PUT", "/:login/:GUID/envcaps/:id", api.UpdateEnvCap},

		// Data Centre Capability Maps

		&rest.Route{"GET", "/#login/:GUID/envcapmaps", api.GetAllEnvCapMaps},

		&rest.Route{"POST", "/:login/:GUID/envcapmaps", api.AddEnvCapMap},

		&rest.Route{"DELETE", "/:login/:GUID/envcapmaps/:id", api.DeleteEnvCapMap},

		&rest.Route{"PUT", "/:login/:GUID/envcapmaps/:id", api.UpdateEnvCapMap},

		// Scripts

		&rest.Route{"GET", "/#login/:GUID/scripts", api.GetAllScripts},

		&rest.Route{"POST", "/:login/:GUID/scripts", api.AddScript},

		&rest.Route{"DELETE", "/:login/:GUID/scripts/:id", api.DeleteScript},

		&rest.Route{"PUT", "/:login/:GUID/scripts/:id", api.UpdateScript},

		// Jobs

		&rest.Route{"GET", "/#login/:GUID/jobs", api.GetAllJobs},

		&rest.Route{"POST", "/#login/:GUID/jobs", api.AddJob},

		&rest.Route{"DELETE", "/#login/:GUID/jobs/kill/:id", api.KillJob},

		&rest.Route{"DELETE", "/#login/:GUID/jobs/:id", api.DeleteJob},

		&rest.Route{"PUT", "/#login/:GUID/jobs/:id", api.UpdateJob},

		// Workers

		&rest.Route{"GET", "/#login/:GUID/workers", api.GetAllWorkers},

		&rest.Route{"POST", "/#login/:GUID/workers", api.AddWorker},

		&rest.Route{"DELETE", "/#login/:GUID/workers/:id", api.DeleteWorker},

		&rest.Route{"PUT", "/#login/:GUID/workers/:id", api.UpdateWorker},

		// JsonObjects

		&rest.Route{"GET", "/#login/:GUID/jsonobjects", api.GetAllJsonObjects},

		&rest.Route{"POST", "/#login/:GUID/jsonobjects", api.AddJsonObject},

		&rest.Route{"DELETE", "/#login/:GUID/jsonobjects/:id", api.DeleteJsonObject},

		&rest.Route{"PUT", "/#login/:GUID/jsonobjects/:id", api.UpdateJsonObject},

		// Repos

		&rest.Route{"GET", "/#login/:GUID/repos", api.GetAllRepos},

		&rest.Route{"POST", "/#login/:GUID/repos", api.AddRepo},

		&rest.Route{"DELETE", "/#login/:GUID/repos/:id", api.DeleteRepo},

		&rest.Route{"PUT", "/#login/:GUID/repos/:id", api.UpdateRepo},

		// Plugins

		&rest.Route{"GET", "/#login/:GUID/plugins", api.GetAllPlugins},

		&rest.Route{"POST", "/#login/:GUID/plugins", api.AddPlugin},

		&rest.Route{"DELETE", "/#login/:GUID/plugins/:id", api.DeletePlugin},

		&rest.Route{"PUT", "/#login/:GUID/plugins/:id", api.UpdatePlugin},

		// Plugin manager

		&rest.Route{"GET", "/#login/:GUID/repoplugins", api.GetAllRepoPlugins},

		&rest.Route{"POST", "/#login/:GUID/repoplugins", api.AddRepoPlugin},

		&rest.Route{"DELETE", "/#login/:GUID/repoplugins/:id", api.DeleteRepoPlugin},

		&rest.Route{"PUT", "/#login/:GUID/repoplugins/:id", api.UpdateRepoPlugin},

		// Files

		&rest.Route{"GET", "/#login/:GUID/files", api.GetAllFiles},

		&rest.Route{"POST", "/#login/:GUID/files", api.AddFile},

		&rest.Route{"DELETE", "/#login/:GUID/files/:id", api.DeleteFile},

		&rest.Route{"PUT", "/#login/:GUID/files/:id", api.UpdateFile},

		// Notifications

		&rest.Route{"GET", "/#login/:GUID/outputlines", api.GetAllOutputLines},

		&rest.Route{"POST", "/#login/:GUID/outputlines", api.AddOutputLine},

		&rest.Route{"DELETE", "/#login/:GUID/outputlines/:id",
			api.DeleteOutputLine},

		/* No changing output logs
		   // //&rest.Route{"PUT", "/#login/:GUID/outputlines/:id", // //
		   // api.UpdateOutputLine},
		*/

		// Generic

		&rest.Route{"GET", "/#login/:GUID/#endpoint",
			api.GenericGetEndpoint},

		&rest.Route{"GET", "/#login/:GUID/#endpoint/#subitem",
			api.GenericGetEndpoint},

		&rest.Route{"POST", "/#login/:GUID/#endpoint",
			api.GenericPostEndpoint},

		&rest.Route{"POST", "/#login/:GUID/#endpoint/#subitem",
			api.GenericPostEndpoint},

		&rest.Route{"DELETE", "/#login/:GUID/#endpoint/:id",
			api.GenericDeleteEndpoint},

		&rest.Route{"DELETE", "/#login/:GUID/#endpoint/#subitem/#id",
			api.GenericDeleteEndpoint},

		&rest.Route{"PUT", "/#login/:GUID/#endpoint", api.GenericPutEndpoint},

		&rest.Route{"PUT", "/#login/:GUID/#endpoint/#subitem",
			api.GenericPutEndpoint},
	)

	// Add the REST API handle
	http.Handle("/api/", http.StripPrefix("/api", &handler))

	// Add the Web Server handle
	changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// HTTP 1.1
			w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
			// HTTP 1.0
			w.Header().Add("Pragma", "no-cache")
			// Proxies
			w.Header().Add("Expires", "0")
			h.ServeHTTP(w, r)
		}
	}
	fs := http.FileServer(http.Dir(config.StaticContent))
	fs = changeHeaderThenServe(fs)
	http.Handle("/manager/", http.StripPrefix("/manager", fs))

	// Add the Manager template handle
	http.HandleFunc("/manager/run", api.serveRunTemplate)

	// Add the Admin template handle
	http.HandleFunc("/manager/admin", api.serveRunTemplate)

	s := http.Server{}
	//s.SetKeepAlivesEnabled(false)
	s.Addr = config.ListenAddr

	if config.SSLEnabled {
		if err := s.ListenAndServeTLS(config.SSLCertFile,
			config.SSLKeyFile); err != nil {
			logit(err.Error())
		}
	} else {
		if err := s.ListenAndServe(); err != nil {
			logit(err.Error())
		}
	}
}
