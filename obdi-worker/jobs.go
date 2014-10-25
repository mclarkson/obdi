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
    "github.com/mclarkson/obdi/external/ant0ine/go-json-rest/rest"
    //"sync"
    "fmt"
    "syscall"
    //"time"
)

func (api *Api) ShowJobs(w rest.ResponseWriter, r *rest.Request) {
    w.WriteJson( api.Jobs() )
}

func (api *Api) AddJob(w rest.ResponseWriter, r *rest.Request) {

    // Decode json post data into JobIn struct

    logit( fmt.Sprintf("Connection from %s",r.RemoteAddr) )

    job := JobIn{}
    if err := r.DecodeJsonPayload(&job); err != nil {
        rest.Error(w, "Invalid data format received.", 400)
        return
    } else if job.JobID == 0 {
        rest.Error(w, "Incorrect data format received.", 400)
        return
    }

    // Check the password matches

    if job.Key != config.WorkerKey {
        rest.Error(w, "Invalid key", 400)
        return
    }

    // Add the job to the job list
    api.AppendJob( job )

    if api.Guid() == "" {
        api.loginmutex.Lock()
        if err := api.Login(); err != nil {
            // Can't send this error to the Manager so must return it here
            logit( fmt.Sprintf("Error: %s", err.Error()) )
            rest.Error(w, err.Error(), 400 )
            api.loginmutex.Unlock()
            return
        }
        api.loginmutex.Unlock()
    }

    if err := api.sendStatus( job, JobOut{
        Status:         STATUS_NOTSTARTED,
        StatusReason:   "About to start job",
        StatusPercent:  0,
        Errors:         0,
    } ); err != nil {
        logit( fmt.Sprintf("Error: %s", err.Error()) )
    }

    w.WriteJson( job )

    //a := fmt.Sprintf("%#v",job)
    //logit(a)

    go api.execCmd( job )
}

func (api *Api) DeleteJob(w rest.ResponseWriter, r *rest.Request) {

    // Decode json post data into JobIn struct

    logit( fmt.Sprintf("Connection from %s",r.RemoteAddr) )

    job := JobIn{}
    if err := r.DecodeJsonPayload(&job); err != nil {
        rest.Error(w, "Invalid data format received.", 400)
        return
    } else if job.JobID == 0 {
        rest.Error(w, "Incorrect data format received.", 400)
        return
    }

    // Check the password matches

    if job.Key != config.WorkerKey {
        rest.Error(w, "Invalid key", 400)
        return
    }

    oldjob, err := api.FindJob( job.JobID )
    if err != nil {
        rest.Error(w, "Job not found", 400)
        return
    }

    // So status can be updated correctly
    api.SetUserCancel( oldjob.JobID )

    // Kill the whole process group (-pid)
    syscall.Kill( int(oldjob.Pid)*-1, syscall.SIGKILL )

    // RemoveJob is done if the Wait fails in execCmd (exec.go)
    // And wait will fail 'cos we just killed it.
    // //api.RemoveJob( oldjob.JobID )

    w.WriteJson( job )
}

