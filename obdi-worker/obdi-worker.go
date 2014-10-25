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

package main

import (
    //"fmt"
    "net/http"
    "github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
)

func main() {
    logit("Worker Starting")

    api := NewApi()

    handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
        //DisableJsonIndent: true,
        //EnableGzip:        true,
    }
    handler.SetRoutes(

        // Show jobs
        rest.RouteObjectMethod("GET", "/jobs", &api, "ShowJobs"),

        // New job
        rest.RouteObjectMethod("POST", "/jobs", &api, "AddJob"),

        // Kill job
        rest.RouteObjectMethod("DELETE", "/jobs", &api, "DeleteJob"),

        // Modify job (?)
        //rest.RouteObjectMethod("PUT", "/jobs", &api, "RunJob"),

    )

    // Add the REST API handle
    http.Handle("/api/", http.StripPrefix("/api", &handler))

    if config.SSLEnabled {
        if err := http.ListenAndServeTLS(config.ListenAddr,
          config.SSLCertFile, config.SSLKeyFile, nil); err != nil {
            logit( err.Error() )
        }
    } else {
        if err := http.ListenAndServe(config.ListenAddr, nil); err != nil {
            logit( err.Error() )
        }
    }
}
