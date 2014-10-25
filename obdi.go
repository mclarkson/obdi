// Deployment Manager - a REST interface and GUI for deploying software
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

    rest.RouteObjectMethod("POST", "/login", &api, "DoLogin"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/logout", &api, "Logout"),

    // ADMIN FUNCTIONS

    // USERS

    rest.RouteObjectMethod("GET", "/:login/:GUID/users",
      &api, "GetAllUsers"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/users",
      &api, "AddUser"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/users/:id",
      &api, "UpdateUser"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/users/:id",
      &api, "DeleteUser"),

    // Data Centres

    rest.RouteObjectMethod("GET", "/:login/:GUID/dcs",
      &api, "GetAllDcs"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/dcs",
      &api, "AddDc"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/dcs/:id",
      &api, "DeleteDc"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/dcs/:id",
      &api, "UpdateDc"),

    // Environments

    rest.RouteObjectMethod("GET", "/#login/:GUID/envs",
      &api, "GetAllEnvs"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/envs",
      &api, "AddEnv"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/envs/:id",
      &api, "DeleteEnv"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/envs/:id",
      &api, "UpdateEnv"),

    // User Permissions

    rest.RouteObjectMethod("GET", "/:login/:GUID/perms",
      &api, "GetAllPerms"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/perms",
      &api, "AddPerm"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/perms/:id",
      &api, "DeletePerm"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/perms/:id",
      &api, "UpdatePerm"),

    // Data Centre Capabilities

    rest.RouteObjectMethod("GET", "/:login/:GUID/dccaps",
      &api, "GetAllDcCaps"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/dccaps",
      &api, "AddDcCap"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/dccaps/:id",
      &api, "DeleteDcCap"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/dccaps/:id",
      &api, "UpdateDcCap"),

    // Data Centre Capability Maps

    rest.RouteObjectMethod("GET", "/:login/:GUID/dccapmaps",
      &api, "GetAllDcCapMaps"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/dccapmaps",
      &api, "AddDcCapMap"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/dccapmaps/:id",
      &api, "DeleteDcCapMap"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/dccapmaps/:id",
      &api, "UpdateDcCapMap"),

    // Environment Capabilities

    rest.RouteObjectMethod("GET", "/:login/:GUID/envcaps",
      &api, "GetAllEnvCaps"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/envcaps",
      &api, "AddEnvCap"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/envcaps/:id",
      &api, "DeleteEnvCap"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/envcaps/:id",
      &api, "UpdateEnvCap"),

    // Data Centre Capability Maps

    rest.RouteObjectMethod("GET", "/:login/:GUID/envcapmaps",
      &api, "GetAllEnvCapMaps"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/envcapmaps",
      &api, "AddEnvCapMap"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/envcapmaps/:id",
      &api, "DeleteEnvCapMap"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/envcapmaps/:id",
      &api, "UpdateEnvCapMap"),

    // Scripts

    rest.RouteObjectMethod("GET", "/#login/:GUID/scripts",
      &api, "GetAllScripts"),

    rest.RouteObjectMethod("POST", "/:login/:GUID/scripts",
      &api, "AddScript"),

    rest.RouteObjectMethod("DELETE", "/:login/:GUID/scripts/:id",
      &api, "DeleteScript"),

    rest.RouteObjectMethod("PUT", "/:login/:GUID/scripts/:id",
      &api, "UpdateScript"),

    // Jobs

    rest.RouteObjectMethod("GET", "/#login/:GUID/jobs",
      &api, "GetAllJobs"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/jobs",
      &api, "AddJob"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/jobs/kill/:id",
      &api, "KillJob"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/jobs/:id",
      &api, "DeleteJob"),

    rest.RouteObjectMethod("PUT", "/#login/:GUID/jobs/:id",
      &api, "UpdateJob"),

    // Plugins

    rest.RouteObjectMethod("GET", "/#login/:GUID/plugins",
      &api, "GetAllPlugins"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/plugins",
      &api, "AddPlugin"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/plugins/:id",
      &api, "DeletePlugin"),

    rest.RouteObjectMethod("PUT", "/#login/:GUID/plugins/:id",
      &api, "UpdatePlugin"),

    // Files

    rest.RouteObjectMethod("GET", "/#login/:GUID/files",
      &api, "GetAllFiles"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/files",
      &api, "AddFile"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/files/:id",
      &api, "DeleteFile"),

    rest.RouteObjectMethod("PUT", "/#login/:GUID/files/:id",
      &api, "UpdateFile"),

    // Notifications

    rest.RouteObjectMethod("GET", "/#login/:GUID/outputlines",
      &api, "GetAllOutputLines"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/outputlines",
      &api, "AddOutputLine"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/outputlines/:id",
      &api, "DeleteOutputLine"),

    /* No changing output logs
       // //rest.RouteObjectMethod("PUT", "/#login/:GUID/outputlines/:id",
       // //    &api, "UpdateOutputLine"),
    */

    // Generic

    rest.RouteObjectMethod("GET", "/#login/:GUID/#endpoint",
      &api, "GenericGetEndpoint"),

    rest.RouteObjectMethod("GET", "/#login/:GUID/#endpoint/#subitem",
      &api, "GenericGetEndpoint"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/#endpoint",
      &api, "GenericPostEndpoint"),

    rest.RouteObjectMethod("POST", "/#login/:GUID/#endpoint/#subitem",
      &api, "GenericPostEndpoint"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/#endpoint",
      &api, "GenericDeleteEndpoint"),

    rest.RouteObjectMethod("DELETE", "/#login/:GUID/#endpoint/#subitem",
      &api, "GenericDeleteEndpoint"),

    rest.RouteObjectMethod("PUT", "/#login/:GUID/#endpoint",
      &api, "GenericPutEndpoint"),

    rest.RouteObjectMethod("PUT", "/#login/:GUID/#endpoint/#subitem",
      &api, "GenericPutEndpoint"),
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
