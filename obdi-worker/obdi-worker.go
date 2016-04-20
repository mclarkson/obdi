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
	//"fmt"
	"crypto/tls"
	"github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
	"net/http"
)

func main() {
	logit("Worker Starting")

	// Assign to global Transport var - accept self-signed certs
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	api := NewApi()

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		//DisableJsonIndent: true,
		EnableGzip: true,
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
