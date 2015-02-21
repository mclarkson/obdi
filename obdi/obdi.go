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

// obdi manages deployment processes
package main

import (
	//"fmt"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}
var apimutex = &sync.Mutex{}

func main() {

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

		&rest.Route{"GET", "/:login/:GUID/envcaps", api.GetAllEnvCaps},

		&rest.Route{"POST", "/:login/:GUID/envcaps", api.AddEnvCap},

		&rest.Route{"DELETE", "/:login/:GUID/envcaps/:id", api.DeleteEnvCap},

		&rest.Route{"PUT", "/:login/:GUID/envcaps/:id", api.UpdateEnvCap},

		// Data Centre Capability Maps

		&rest.Route{"GET", "/:login/:GUID/envcapmaps", api.GetAllEnvCapMaps},

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

	if config.SSLEnabled {
		if err := http.ListenAndServeTLS(config.ListenAddr,
			config.SSLCertFile, config.SSLKeyFile, nil); err != nil {
			logit(err.Error())
		}
	} else {
		if err := http.ListenAndServe(config.ListenAddr, nil); err != nil {
			logit(err.Error())
		}
	}
}
