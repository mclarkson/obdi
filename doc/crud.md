Obdi Core REST Usage
====================

Log In
------

Log in is the admin user:

```
ipport="127.0.0.1:443"

guid=`curl -s \
      -d '{"Login":"admin","Password":"admin"}' \
      http://$ipport/api/login | grep -o "[a-z0-9][^\"]*"`
```

Check that the GUID token was received:

    echo $guid

Log Out
-------

Logging out will invalidate the GUID. Log admin out:

    curl -sk -X POST "http://$ipport/api/admin/$guid/logout"

Data Centres
------------

```
# Add a new data centre (HTTP POST)

curl -s -d '{
    "SysName":"NEWDC",
    "DispName":"New DC"
}' "http://$ipport/api/admin/$guid/dcs"

# View details for all data centres (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dcs"

# View details for a single data centre (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dcs?sys_name=NEWDC"

# Update data centre details (HTTP PUT)

curl -s -X PUT -d '{
    "DispName":"A new dc decsription"}' "http://$ipport/api/admin/$guid/dcs/4"

# Delete a data centre (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/dcs/4"

```

Data Centre Capabilities
------------------------

```
# Add a new capability (HTTP POST)

curl -s -d '{
    "Code":"HAS_RH_REPO_CLONE",
    "Desc":"Has Red Hat repo clone"
}' "http://$ipport/api/admin/$guid/dccaps"

# View details for all capabilities (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dccaps"

# View details for a single capability (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dccaps?code=HAS_RH_REPO_CLONE"

# Update capability details (HTTP PUT)

curl -s -X PUT -d '{
    "Desc":"A new dccap description"}' "http://$ipport/api/admin/$guid/dccaps/1"

# Delete a capability (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/dccaps/1"

```

Environments
------------

```
# Add a new environment (HTTP POST)

curl -s -d '{
    "SysName":"NEWENV",
    "DispName":"New Environment"
    "DcId":1
}' "http://$ipport/api/admin/$guid/envs"

# View details for all environments (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envs"

# View details for a single environment (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envs?sys_name=NEWENV"

# Update environment details (HTTP PUT)

curl -s -X PUT -d '{
    "DispName":"A new env description"}' "http://$ipport/api/admin/$guid/envs/21"

# Delete an environment (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/envs/21"

```

Environment Capabilities
------------------------

```
# Add a new capability (HTTP POST)

curl -s -d '{
    "Code":"HAS_CUSTOM_RPM_REPO",
    "Desc":"Has a custom RPM repository."
    "IsWorkerDef":false
}' "http://$ipport/api/admin/$guid/envcaps"

# View details for all capabilities (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envcaps"

# View details for a single capability (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envcaps?code=HAS_CUSTOM_RPM_REPO"

# Update capability details (HTTP PUT)

curl -s -X PUT -d '{
    "Desc":"A new dccap description"}' "http://$ipport/api/admin/$guid/envcaps/1"

# Delete a capability (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/envcaps/1"

```

Workers
-------

If an environment capability has `IsWorkerDef` set to `true` then a Worker
entry is expected for that capability.

```
# Add a new worker (HTTP POST)

curl -s -d '{
    "EnvId":1,
    "EnvCapId":3,
    "WorkerUrl":"https://ServerNameOrIP:4443/",
    "WorkerKey":"PassW0rd"
}' "http://$ipport/api/admin/$guid/workers"

# View details for all workers (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/workers"

# View details for a single worker (HTTP GET)
# Supports using env_id *and* env_cap_id

curl -s "http://$ipport/api/admin/$guid/workers?env_id=1&env_cap_id=3"

# Update worker details (HTTP PUT)

curl -s -X PUT -d '{
    "WorkerKey":"lOcAlH0St"}' "http://$ipport/api/admin/$guid/workers/1"

# Delete a worker (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/workers/1"

```

User Permissions
----------------

```
# Add a new environment permission (HTTP POST)

curl -s -d '{
    "UserId":"NEWENV",
    "EnvId":"New Permission"
    "Enabled":true
    "Writeable":true
}' "http://$ipport/api/admin/$guid/perms"

# View details for all environments (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/perms"

# View permission details for a single user (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/perms?user_id=2"

# Update permission details (HTTP PUT)

curl -s -X PUT -d '{
    "Enabled":true,"Writeable":true}' "http://$ipport/api/admin/$guid/perms/2"

# Delete a permission (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/perms/2"

```

Data Centre Capability Maps
---------------------------

```
# Add a new mapping (HTTP POST)

curl -s -d '{
    "DcId":'"1"',
    "DcCapId":'"1"'
}' "http://$ipport/api/admin/$guid/dccapmaps"

# View details for all mappings (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dccapmaps"

# View permission details for a single data centre (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/dccapmaps?dc_id=1"

# Update a mapping (HTTP PUT)

curl -s -X PUT -d '{
    "DcCapId":'"1"}' "http://$ipport/api/admin/$guid/dccapmaps/1"

# Delete a mapping (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/dccapmaps/1"

```

Environment Capability Maps
---------------------------

```
# Add a new mapping (HTTP POST)

curl -s -d '{
    "EnvId":'"1"',
    "EnvCapId":'"1"'
}' "http://$ipport/api/admin/$guid/envcapmaps"

# View details for all mappings (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envcapmaps"

# View permission details for a single data centre (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/envcapmaps?env_id=1"

# Update a mapping (HTTP PUT)

curl -s -X PUT -d '{
    "DcCapId":'"1"}' "http://$ipport/api/admin/$guid/envcapmaps/1"

# Delete a mapping (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/envcapmaps/1"

```

Jobs
----

```
# Add a new job (HTTP POST)

curl -s -d '{
    "ScriptId":1,
    "Args":"-s 4",
    "EnvVars":"A=1 B=2 C=\"Hello there\" D=4"
}' "http://$ipport/api/a.user/$guid/jobs"

# View details for all jobs (HTTP GET)

curl -s "http://$ipport/api/a.user/$guid/jobs"

# View permission details for a single job (HTTP GET)

curl -s "http://$ipport/api/a.user/$guid/jobs?job_id=1"

# Update a job (HTTP PUT)

curl -s -X PUT -d '{
    "Status":1}' "http://$ipport/api/a.user/$guid/jobs/1"

# Kill a running job

curl -i -X DELETE "http://$ipport/api/a.user/$guid/jobs/kill/1"

# Delete a job (HTTP DELETE) This does not kill a running job!

curl -i -X DELETE "http://$ipport/api/a.user/$guid/jobs/1"

```

Scripts
-------


```
# Add a new scripts (HTTP POST)
# Source is a script, or binary, encoded in base64

curl -s -d '{
    "Name":"env.sh",
    "Desc":"Output environment variables",
    "Source":"IyEvYmluL2Jhc2gKClBBVEg9JFBBVEg6L2JpbjovdXNyL2JpbgoKZWNobyAiSGkiCgpzbGVlcCAxCgplY2hvICRACgpzbGVlcCAxCgplbnYKCnNsZWVwIDEKCmVjaG8gIkJ5ZSIKCmV4aXQgMAoK"
}' "http://$ipport/api/admin/$guid/scripts"

# View details for all scripts (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/scripts"

# View permission details for a single script (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/scripts?script_id=1"

# Same as previous but exclude the 'Source' column

curl -s "http://$ipport/api/admin/$guid/scripts?script_id=1&nosource=1"

# Update a script (HTTP PUT)

curl -s -X PUT -d '{"Name":"env_old.sh"}' "http://$ipport/api/admin/$guid/scripts/1"

# Delete a script (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/scripts/1"

```

Output Lines
------------

Command line output from scripts. Output lines are stored over multiple rows.
Output lines can be update and deleted, but they cannot be updated (changed).

```
# Add a new outputline (HTTP POST)

curl -s -d '{
    "Serial":1,
    "JobID":1,
    "Text":"Hello"
}' "http://$ipport/api/a.user/$guid/outputlines"

# View details for all outputlines (HTTP GET)

curl -s "http://$ipport/api/a.user/$guid/outputlines"

# View outputlines for a single job (HTTP GET)
# Add '&top=5' to see the top 5 lines or
# '&bottom=5' to see the bottom 5 lines.

curl -s "http://$ipport/api/a.user/$guid/outputlines?job_id=1"

# Update an outputline (HTTP PUT)

No method!

# Delete all outputlines for a job ID (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/a.user/$guid/scripts/1"
                                                           ^
                            This is a job_id, not an outputline_id

```

Plugins
-------

```
# Add a new plugin (HTTP POST)

curl -s -d '{
    "Name":"mysuperplugin",
    "Desc":"My super plugin",
    "Parent":0,
    "SbName":"Super Thing"
}' "http://$ipport/api/admin/$guid/plugins"

# View details for all plugins (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/plugins"

# View permission details for a single script (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/plugins?name=mysuperplugin"

# Update a script (HTTP PUT)

curl -s -X PUT -d '{"Desc":"The super plugin"}'
    "http://$ipport/api/admin/$guid/plugins/1"

# Delete a plugin (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/plugins/1"

```

Files
-----

Files that belong to plugins.

```
# Add a new scripts (HTTP POST)
# Source is a script or binary encoded in base64

curl -s -d '{
    "Name":"afile.js",
    "Desc":"a test file",
    "Type":1,
    "PluginID":1,
    "Url":"mysuperplugin/js/controller/afile.js"
}' "http://$ipport/api/admin/$guid/files"

# View details for all scripts (HTTP GET)

curl -s "http://$ipport/api/admin/$guid/files"

# View permission details for a single script name (HTTP GET)
# This might return many scripts with the same name for different plugins

curl -s "http://$ipport/api/admin/$guid/files?name=afile.js"

# Limit the previous file search to a specific plugin id

curl -s "http://$ipport/api/admin/$guid/files?name=afile.js&plugin_id=1"

# Update a script (HTTP PUT)

curl -s -X PUT -d '{"Name":"a_file.js_old"}' "http://$ipport/api/admin/$guid/files/1"

# Delete a script (HTTP DELETE)

curl -i -X DELETE "http://$ipport/api/admin/$guid/files/1"

```

Worker REST Interface
---------------------

The Manager submits jobs to the Worker.

```
# Add a new job (HTTP POST)
# ScriptSource is a base64 encoded script/binary.

curl -s -d '{
    "ScriptSource":"IyEvYmluL2Jhc2gKClBBVEg9JFBBVEg6L2JpbjovdXNyL2JpbgoKZWNobyAiSGkiCgpzbGVlcCAxCgplY2hvICRACgpzbGVlcCAxCgplbnYKCnNsZWVwIDEKCmVjaG8gIkJ5ZSIKCmV4aXQgMAoK",
    "EnvVars":"A=1 B=\"Hi There\" C=3",
    "JobID":2,
    "Key":"serverkey",
    "Args":"-a \"hi there\" -x 1"
}' "http://$ipport/api/a.user/$guid/jobs"

# View details for all jobs (HTTP GET)

curl -s "http://$ipport/api/a.user/$guid/jobs"

# View permission details for a single job (HTTP GET)

curl -s "http://$ipport/api/a.user/$guid/jobs?job_id=1"

# Update a job (HTTP PUT)

There is no Update enpoint!

# Delete a job (HTTP DELETE). This kills all processes in the process group for
# this job.

curl -i -X DELETE '{
    "JobID":1,
    "Key":"serverkey"
}' "http://$ipport/api/a.user/$guid/jobs/1"

```

