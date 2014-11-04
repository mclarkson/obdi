#!/bin/bash

proto="https"
opts="-k -s" # don't check ssl cert, silent
ipport="127.0.0.1:443"
guid=`curl $opts -d '{"Login":"admin","Password":"admin"}' \
    $proto://$ipport/api/login | grep -o "[a-z0-9][^\"]*"`

echo "GUID=$guid"

# Create users
curl $opts -d '{
    "login":"mark.clarkson",
    "passHash":"password",
    "forename":"Mark",
    "surname":"Clarkson",
    "email":"mark.clarkson@nowhere",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

curl $opts -d '{
    "login":"adrian.smart",
    "passHash":"password",
    "forename":"Adrian",
    "surname":"Smart",
    "email":"adrian.smart@nowhere",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

curl $opts -d '{
    "login":"marcelo.black",
    "passHash":"password",
    "forename":"Marcelo",
    "surname":"Black",
    "email":"marcelo.black@nowhere",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

curl $opts -d '{
    "login":"patrick.stewart",
    "passHash":"password",
    "forename":"Patrick",
    "surname":"Stewart",
    "email":"patrick.stewart@nowhere",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

curl $opts -d '{
    "login":"hasham.amla",
    "passHash":"password",
    "forename":"Hasham",
    "surname":"Amla",
    "email":"hasham.amla@nowhere",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

curl $opts -d '{
    "login":"worker",
    "passHash":"p1ssword",
    "forename":"Remote",
    "surname":"Worker Daemon",
    "email":"no@email",
    "enabled":true}' "$proto://$ipport/api/admin/$guid/users"

# Create Data Centres

curl $opts -d '{
    "SysName":"OFFICE",
    "DispName":"Somewhere Street Office"
    }' "$proto://$ipport/api/admin/$guid/dcs"

curl $opts -d '{
    "SysName":"PROVIDER1",
    "DispName":"Cloud Provider in Docklands"
    }' "$proto://$ipport/api/admin/$guid/dcs"

# Create Environments

## Add DEV TEST STAGE and PROD to OFFICE dc

dcid=`curl $opts "$proto://$ipport/api/admin/$guid/dcs?sys_name=OFFICE" | grep Id | grep -o "[0-9]"`

curl $opts -d '{
    "SysName":"DEV",
    "DispName":"Development",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"TEST",
    "DispName":"Test",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"STAGE",
    "DispName":"Stage",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"PROD",
    "DispName":"Production",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

## Add DEV TEST STAGE and PROD to PROVIDER1 dc

dcid=`curl $opts "$proto://$ipport/api/admin/$guid/dcs?sys_name=PROVIDER1" | grep Id | grep -o "[0-9]"`

curl $opts -d '{
    "SysName":"DEV",
    "DispName":"Development",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"TEST",
    "DispName":"Test",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"STAGE",
    "DispName":"Stage",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

curl $opts -d '{
    "SysName":"PROD",
    "DispName":"Production",
    "DcId":'"$dcid"',
    "WorkerUrl":"https://127.0.0.1:4443/",
    "WorkerKey":"password"
}' "$proto://$ipport/api/admin/$guid/envs"

# Edit Perms

# ---------------------------------------------------------------------------
update_perm()
# ---------------------------------------------------------------------------
# $1 - Perm ID (Integer)
# $2 - Enabled (true|false)
# $3 - Writeable (true|false)
{
    echo curl $opts -X PUT -d '{
        "Enabled":'"$2"',
        "Writeable":'"$3"'
    }' "$proto://$ipport/api/admin/$guid/perms/$1"

    curl $opts -X PUT -d '{
        "Enabled":'"$2"',
        "Writeable":'"$3"'
    }' "$proto://$ipport/api/admin/$guid/perms/$1"
}

# ---------------------------------------------------------------------------
add_perm()
# ---------------------------------------------------------------------------
# $1 - login (text)
# $2 - data centre (text)
# $3 - environment (text)
# $4 - Enabled (true|false)
# $5 - Writeable (true|false)
{
    userid=`curl $opts "$proto://$ipport/api/admin/$guid/users?login=$1" | grep Id | grep -o "[0-9]"`
    dcid=`curl $opts "$proto://$ipport/api/admin/$guid/dcs?sys_name=$2" | grep Id | grep -o "[0-9]"`
    envid=`curl $opts "$proto://$ipport/api/admin/$guid/envs?sys_name=$3&dc_id=$dcid" | grep -w Id | grep -o "[0-9]"`

    echo
    echo "$1 $2 $3 $4 $5"
    echo curl $opts -d '{
        "UserId":'"$userid"',
        "EnvId":'"$envid"',
        "Enabled":true,
        "Writeable":true
    }' "$proto://$ipport/api/admin/$guid/perms"

    curl $opts -d '{
        "UserId":'"$userid"',
        "EnvId":'"$envid"',
        "Enabled":'"$4"',
        "Writeable":'"$5"'
    }' "$proto://$ipport/api/admin/$guid/perms"
}

#add_perm adrian.smart OFFICE DEV true true
add_perm adrian.smart OFFICE PROD true true
add_perm adrian.smart OFFICE STAGE true true
add_perm adrian.smart OFFICE TEST true true
add_perm adrian.smart PROVIDER1 DEV true false
add_perm adrian.smart PROVIDER1 PROD true false
add_perm adrian.smart PROVIDER1 STAGE true false
add_perm adrian.smart PROVIDER1 TEST true false

add_perm hasham.amla OFFICE DEV true true
#add_perm hasham.amla OFFICE PROD true true
add_perm hasham.amla OFFICE STAGE true true
add_perm hasham.amla OFFICE TEST true true
add_perm hasham.amla PROVIDER1 DEV true false
add_perm hasham.amla PROVIDER1 PROD true false
add_perm hasham.amla PROVIDER1 STAGE true false
add_perm hasham.amla PROVIDER1 TEST true false

add_perm marcelo.black OFFICE DEV true true
add_perm marcelo.black OFFICE PROD true true
#add_perm marcelo.black OFFICE STAGE true true
add_perm marcelo.black OFFICE TEST true true
add_perm marcelo.black PROVIDER1 DEV true false
add_perm marcelo.black PROVIDER1 PROD true false
add_perm marcelo.black PROVIDER1 STAGE true false
add_perm marcelo.black PROVIDER1 TEST true false

add_perm mark.clarkson OFFICE DEV true true
add_perm mark.clarkson OFFICE PROD true true
add_perm mark.clarkson OFFICE STAGE true true
#add_perm mark.clarkson OFFICE TEST true true
add_perm mark.clarkson PROVIDER1 DEV true false
add_perm mark.clarkson PROVIDER1 PROD true false
add_perm mark.clarkson PROVIDER1 STAGE true false
add_perm mark.clarkson PROVIDER1 TEST true false

add_perm patrick.stewart OFFICE DEV true true
add_perm patrick.stewart OFFICE PROD true true
add_perm patrick.stewart OFFICE STAGE true true
add_perm patrick.stewart OFFICE TEST true true
#add_perm patrick.stewart PROVIDER1 DEV true false
add_perm patrick.stewart PROVIDER1 PROD true false
add_perm patrick.stewart PROVIDER1 STAGE true false
add_perm patrick.stewart PROVIDER1 TEST true false

#update_perm 5 false false
#update_perm 5 true false

curl $opts -d '{ "Code":"HAS_RH_REPO_CLONE","Desc":"The data centre CLONES a Red Hat YUM repository from the provider." }' "$proto://$ipport/api/admin/$guid/dccaps"
curl $opts -d '{ "Code":"HAS_CENTOS_REPO_CLONE","Desc":"The data centre has a CentOS YUM repository." }' "$proto://$ipport/api/admin/$guid/dccaps"
curl $opts -d '{ "Code":"HAS_HOSTED_RH_REPO","Desc":"The data centre has a Red Hat YUM repository hosted by the provider." }' "$proto://$ipport/api/admin/$guid/dccaps"

curl $opts -d '{ "Code":"HAS_YUM_REPO","Desc":"Has a small YUM repo for custom RPMs" }' "$proto://$ipport/api/admin/$guid/envcaps"
curl $opts -d '{ "Code":"HAS_CUSTOM_REPO","Desc":"Has a custom repo created from Jenkins" }' "$proto://$ipport/api/admin/$guid/envcaps"
curl $opts -d '{ "Code":"HAS_SALT_REPO","Desc":"Has a Salt repo cloned from Beanstalk" }' "$proto://$ipport/api/admin/$guid/envcaps"
curl $opts -d '{ "Code":"HAS_SERVER_STATUS_LIST","Desc":"Has a server status list view" }' "$proto://$ipport/api/admin/$guid/envcaps"


# ---------------------------------------------------------------------------
add_dccapmap()
# ---------------------------------------------------------------------------
# $1 - data centre (text)
# $2 - data centre capability (text)
{
    dcid=`curl $opts "$proto://$ipport/api/admin/$guid/dcs?sys_name=$1" | grep Id | grep -o "[0-9]"`
    dccapid=`curl $opts "$proto://$ipport/api/admin/$guid/dccaps?code=$2" | grep Id | grep -o "[0-9]"`

    echo
    echo curl $opts -d '{
        "DcId":'"$dcid"',
        "DcCapId":'"$dccapid"',
    }' "$proto://$ipport/api/admin/$guid/dccapmaps"

    curl $opts -d '{
        "DcId":'"$dcid"',
        "DcCapId":'"$dccapid"'
    }' "$proto://$ipport/api/admin/$guid/dccapmaps"
}

add_dccapmap PROVIDER1 HAS_HOSTED_RH_REPO
add_dccapmap OFFICE HAS_RH_REPO_CLONE

# ---------------------------------------------------------------------------
add_envcapmap()
# ---------------------------------------------------------------------------
# $1 - environment (text)
# $2 - data centre (text)
# $3 - environment capability (text)
{
    dcid=`curl $opts "$proto://$ipport/api/admin/$guid/dcs?sys_name=$2" | grep Id | grep -o "[0-9]"`
    envid=`curl $opts "$proto://$ipport/api/admin/$guid/envs?sys_name=$1&dc_id=$dcid" | grep -w Id | grep -o "[0-9]"`
    envcapid=`curl $opts "$proto://$ipport/api/admin/$guid/envcaps?code=$3" | grep -w Id | grep -o "[0-9]"`

    echo
    echo curl $opts -d '{
        "EnvId":'"$envid"',
        "EnvCapId":'"$envcapid"',
    }' "$proto://$ipport/api/admin/$guid/envcapmaps"

    curl $opts -d '{
        "EnvId":'"$envid"',
        "EnvCapId":'"$envcapid"'
    }' "$proto://$ipport/api/admin/$guid/envcapmaps"
}

add_envcapmap DEV OFFICE HAS_YUM_REPO
add_envcapmap DEV OFFICE HAS_CUSTOM_REPO
add_envcapmap DEV OFFICE HAS_SALT_REPO
add_envcapmap PROD OFFICE HAS_YUM_REPO
add_envcapmap PROD OFFICE HAS_CUSTOM_REPO
add_envcapmap PROD OFFICE HAS_SALT_REPO
add_envcapmap STAGE OFFICE HAS_YUM_REPO
add_envcapmap STAGE OFFICE HAS_CUSTOM_REPO
add_envcapmap STAGE OFFICE HAS_SALT_REPO
add_envcapmap TEST OFFICE HAS_YUM_REPO
add_envcapmap TEST OFFICE HAS_CUSTOM_REPO
add_envcapmap TEST OFFICE HAS_SALT_REPO

add_envcapmap DEV PROVIDER1 HAS_YUM_REPO
add_envcapmap DEV PROVIDER1 HAS_CUSTOM_REPO
add_envcapmap DEV PROVIDER1 HAS_SALT_REPO
add_envcapmap PROD PROVIDER1 HAS_YUM_REPO
add_envcapmap PROD PROVIDER1 HAS_CUSTOM_REPO
add_envcapmap PROD PROVIDER1 HAS_SALT_REPO
add_envcapmap STAGE PROVIDER1 HAS_YUM_REPO
add_envcapmap STAGE PROVIDER1 HAS_CUSTOM_REPO
add_envcapmap STAGE PROVIDER1 HAS_SALT_REPO
add_envcapmap TEST PROVIDER1 HAS_YUM_REPO
add_envcapmap TEST PROVIDER1 HAS_CUSTOM_REPO
add_envcapmap TEST PROVIDER1 HAS_SALT_REPO

# Scripts (for testing)
# 
# curl $opts -d '{
# "Name":"env.sh",
# "Desc":"Output environment variables",
# "Source":"IyEvYmluL2Jhc2gKClBBVEg9JFBBVEg6L2JpbjovdXNyL2JpbgoKZWNobyAiSGkiCgpzbGVlcCAxCgplY2hvICRACgpzbGVlcCAxCgplbnYKCnNsZWVwIDEKCmVjaG8gIkJ5ZSIKCmV4aXQgMAoK"
# }' $proto://$ipport/api/admin/$guid/scripts
# 
# curl $opts -d '{                                                                                                 
# "Name":"bogus.sh",
# "Desc":"A bogus script",
# "Source":"IyEvYmluL2Jhc2gKClBB"                                                                                                                                
# }' $proto://$ipport/api/admin/$guid/scripts
# 
# curl $opts -d '{
# "Name":"env_long.sh",
# "Desc":"Env output that takes around 300 seconds",
# "Source":"IyEvYmluL2Jhc2gKClBBVEg9JFBBVEg6L2JpbjovdXNyL2JpbgoKZWNobyAiSGkiCgpzbGVlcCA1CgplY2hvICRACgpzbGVlcCA0CgplbnYKCnNsZWVwIDI5MQoKZWNobyAiQnllIgoKZXhpdCAwCgo="
# }' $proto://$ipport/api/admin/$guid/scripts
# 
# curl $opts -d '{
# "Name":"run_long.sh",
# "Desc":"A script that spawns a child process that runs forever",
# "Source":"IyEvYmluL2Jhc2gKCiMgU3Bhd24gYSBjaGlsZCBwcm9jZXNzLiBUaGlzIHNob3VsZCBiZSBraWxsZWQgdG9vLgoKKCBhPTEKICB3aGlsZSB0cnVlOyBkbwogICAgICBlY2hvICQoKGErKykpIAogICAgICBzbGVlcCAxCiAgZG9uZQopICYKCndhaXQKCiMgVGhpcyBwcm9jZXNzIHdpbGwgYmUga2lsbGVkIHNvIHdlJ2xsIG5ldmVyIGdldCBoZXJlLgoKZWNobyAiRG9uZSEiCgpleGl0IDAK"
# }' $proto://$ipport/api/admin/$guid/scripts

# Plugins

# Dashboard

curl -k -d '{
    "Name":"dashboard",
    "Desc":"A Dashboard framework",
    "HasView":1,
    "Parent":""
}' $proto://$ipport/api/admin/$guid/plugins | tee /tmp/fshodifha
# Grab the id of the last insert
id=`grep Id /tmp/fshodifha | grep -Eo "[0-9]+"`
# Add the AJS controller file
curl -k -d '{
    "Name":"dashboard.js",
    "Desc":"Dashboard controller",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"dashboard/js/controllers/dashboard.js"
}' $proto://$ipport/api/admin/$guid/files

# dash-jobs

curl -k -d '{
    "Name":"dash-jobs",
    "Desc":"A Dashboard status plugin for Jobs",
    "HasView":1,
    "Parent":"dashboard"
}' $proto://$ipport/api/admin/$guid/plugins

# # jobs
# 
# curl -k -d '{
#     "Name":"jobs",
#     "Desc":"Job management",
#     "HasView":1,
#     "Parent":""
# }' $proto://$ipport/api/admin/$guid/plugins | tee /tmp/fshodifha
# # Grab the id of the last insert
# id=`grep Id /tmp/fshodifha | grep -Eo "[0-9]+"`
# # Add the AJS controller file
# curl -k -d '{
#     "Name":"jobs.js",
#     "Desc":"Manages the job queue",
#     "Type":1,
#     "PluginId":'"$id"',
#     "Url":"jobs/js/controllers/jobs.js"
# }' $proto://$ipport/api/admin/$guid/files
# curl -k -d '{
#     "Name":"outputlines.js",
#     "Desc":"Shows job log",
#     "Type":1,
#     "PluginId":'"$id"',
#     "Url":"jobs/js/controllers/outputlines.js"
# }' $proto://$ipport/api/admin/$guid/files

# salt stub

# curl -k -d '{
#     "Name":"salt",
#     "Desc":"Salt main sidebar item",
#     "HasView":2,
#     "Parent":""
# }' $proto://$ipport/api/admin/$guid/plugins

# salt update

curl -k -d '{
    "Name":"saltupdate",
    "Desc":"Update the Salt GIT sources",
    "HasView":1,
    "Parent":"salt"
}' $proto://$ipport/api/admin/$guid/plugins

# # salt config server
# 
# curl -k -d '{
#     "Name":"saltconfigserver",
#     "Desc":"Salt New Server plugin",
#     "HasView":1,
#     "Parent":"salt"
# }' $proto://$ipport/api/admin/$guid/plugins | tee /tmp/fshodifha
# # Grab the id of the last insert
# id=`grep Id /tmp/fshodifha | grep -Eo "[0-9]+"`
# # Add the AJS controller file
# curl -k -d '{
#     "Name":"saltconfigserver.js",
#     "Desc":"Controller for Salt New Server",
#     "Type":1,
#     "PluginId":'"$id"',
#     "Url":"saltconfigserver/js/controllers/saltconfigserver.js"
# }' $proto://$ipport/api/admin/$guid/files
# # Add the scripts
# curl -k -d '{
#     "Desc": "Return the grains for a server from the Salt Cache. Arg1 - Salt ID",
#     "Name": "salt-grains-cache.sh",
#     "Source": "IyEvYmluL2Jhc2gKCmV4cG9ydCBQQVRIPS9iaW46L3Vzci9iaW4KCmlmICEgd2hpY2ggc2FsdCA+JiAvZGV2L251bGw7IHRoZW4KICAgIGVjaG8gJ3siRXJyb3IiOiJTYWx0IGJpbmFyeSBub3QgZm91bmQuIn0nCiAgICBleGl0IDEKZmkKCnNhbHRfaWQ9JDEKCiMgVGhlIHNlcnZlciBpcyBvZmZsaW5lLCBzbyBkb24ndDoKIyAjb3V0cHV0PWBzYWx0IC0tb3V0cHV0PWpzb24gIiRzYWx0X2lkIiBncmFpbnMuaXRlbXNgCm91dHB1dD0KCltbIC16ICRvdXRwdXQgXV0gJiYgewoKb3V0cHV0PWBweXRob24gPCA8KGNhdCA8PEVuRAojIS91c3IvYmluL2VudiBweXRob24KaW1wb3J0IGpzb24KaW1wb3J0IG9zCmltcG9ydCBzeXMKaW1wb3J0IGdldG9wdAoKaW1wb3J0IHNhbHQuY29uZmlnCmltcG9ydCBzYWx0LnJ1bm5lcgoKaWYgX19uYW1lX18gPT0gJ19fbWFpbl9fJzoKICAgIF9fb3B0c19fID0gc2FsdC5jb25maWcubWFzdGVyX2NvbmZpZygKICAgICAgICAgICAgb3MuZW52aXJvbi5nZXQoJ1NBTFRfTUFTVEVSX0NPTkZJRycsICcvZXRjL3NhbHQvbWluaW9uJykpCiAgICBydW5uZXIgPSBzYWx0LnJ1bm5lci5SdW5uZXIoX19vcHRzX18pCgogICAgc3Rkb3V0X2JhayA9IHN5cy5zdGRvdXQKICAgIHdpdGggb3Blbihvcy5kZXZudWxsLCAnd2InKSBhcyBmOgogICAgICAgIHN5cy5zdGRvdXQgPSBmCiAgICAgICAgaXRlbXMgPSBydW5uZXIuY21kKCJjYWNoZS5ncmFpbnMiLCBbIiRzYWx0X2lkIl0pCiAgICBzeXMuc3Rkb3V0ID0gc3Rkb3V0X2JhawoKICAgIHByaW50IGpzb24uZHVtcHMoaXRlbXMpCkVuRAopYAp9CgplY2hvICRvdXRwdXQK"
# }' $proto://$ipport/api/admin/$guid/scripts
# curl -k -d '{
#     "Desc": "Return the grains for a server. Arg1 - Salt ID",
#     "Name": "salt-grains.sh",
#     "Source": "IyEvYmluL2Jhc2gKCmV4cG9ydCBQQVRIPS9iaW46L3Vzci9iaW4KCmlmICEgd2hpY2ggc2FsdCA+JiAvZGV2L251bGw7IHRoZW4KICAgIGVjaG8gJ3siRXJyb3IiOiJTYWx0IGJpbmFyeSBub3QgZm91bmQuIn0nCiAgICBleGl0IDEKZmkKCnNhbHRfaWQ9JDEKCm91dHB1dD1gc2FsdCAtLW91dHB1dD1qc29uICIkc2FsdF9pZCIgZ3JhaW5zLml0ZW1zYAoKW1sgLXogJG91dHB1dCBdXSAmJiB7CgpvdXRwdXQ9YHB5dGhvbiA8IDwoY2F0IDw8RW5ECiMhL3Vzci9iaW4vZW52IHB5dGhvbgppbXBvcnQganNvbgppbXBvcnQgb3MKaW1wb3J0IHN5cwppbXBvcnQgZ2V0b3B0CgppbXBvcnQgc2FsdC5jb25maWcKaW1wb3J0IHNhbHQucnVubmVyCgppZiBfX25hbWVfXyA9PSAnX19tYWluX18nOgogICAgX19vcHRzX18gPSBzYWx0LmNvbmZpZy5tYXN0ZXJfY29uZmlnKAogICAgICAgICAgICBvcy5lbnZpcm9uLmdldCgnU0FMVF9NQVNURVJfQ09ORklHJywgJy9ldGMvc2FsdC9taW5pb24nKSkKICAgIHJ1bm5lciA9IHNhbHQucnVubmVyLlJ1bm5lcihfX29wdHNfXykKCiAgICBzdGRvdXRfYmFrID0gc3lzLnN0ZG91dAogICAgd2l0aCBvcGVuKG9zLmRldm51bGwsICd3YicpIGFzIGY6CiAgICAgICAgc3lzLnN0ZG91dCA9IGYKICAgICAgICBpdGVtcyA9IHJ1bm5lci5jbWQoImNhY2hlLmdyYWlucyIsIFsiJHNhbHRfaWQiXSkKICAgIHN5cy5zdGRvdXQgPSBzdGRvdXRfYmFrCgogICAgcHJpbnQganNvbi5kdW1wcyhpdGVtcykKRW5ECilgCn0KCmVjaG8gJG91dHB1dAo="
# }' $proto://$ipport/api/admin/$guid/scripts
# curl -k -d '{
#     "Desc": "Returns a list of servers, dead or alive. Arg 1 - dc, Arg 2 - env",
#     "Name": "serverlist.py",
#     "Source": "IyEvdXNyL2Jpbi9lbnYgcHl0aG9uCiMKIyBEZXBsb3ltZW50IE1hbmFnZXIgLSBhIFJFU1QgaW50ZXJmYWNlIGFuZCBHVUkgZm9yIGRlcGxveWluZyBzb2Z0d2FyZQojIENvcHlyaWdodCAoQykgMjAxNCAgTWFyayBDbGFya3NvbgojCiMgVGhpcyBwcm9ncmFtIGlzIGZyZWUgc29mdHdhcmU6IHlvdSBjYW4gcmVkaXN0cmlidXRlIGl0IGFuZC9vciBtb2RpZnkKIyBpdCB1bmRlciB0aGUgdGVybXMgb2YgdGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGFzIHB1Ymxpc2hlZCBieQojIHRoZSBGcmVlIFNvZnR3YXJlIEZvdW5kYXRpb24sIGVpdGhlciB2ZXJzaW9uIDMgb2YgdGhlIExpY2Vuc2UsIG9yCiMgKGF0IHlvdXIgb3B0aW9uKSBhbnkgbGF0ZXIgdmVyc2lvbi4KIwojIFRoaXMgcHJvZ3JhbSBpcyBkaXN0cmlidXRlZCBpbiB0aGUgaG9wZSB0aGF0IGl0IHdpbGwgYmUgdXNlZnVsLAojIGJ1dCBXSVRIT1VUIEFOWSBXQVJSQU5UWTsgd2l0aG91dCBldmVuIHRoZSBpbXBsaWVkIHdhcnJhbnR5IG9mCiMgTUVSQ0hBTlRBQklMSVRZIG9yIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFLiAgU2VlIHRoZQojIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGZvciBtb3JlIGRldGFpbHMuCiMKIyBZb3Ugc2hvdWxkIGhhdmUgcmVjZWl2ZWQgYSBjb3B5IG9mIHRoZSBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZQojIGFsb25nIHdpdGggdGhpcyBwcm9ncmFtLiAgSWYgbm90LCBzZWUgPGh0dHA6Ly93d3cuZ251Lm9yZy9saWNlbnNlcy8+LgoKaW1wb3J0IGpzb24KaW1wb3J0IG9zCmltcG9ydCBzeXMKaW1wb3J0IGdldG9wdAppbXBvcnQgcmUKCmltcG9ydCBzYWx0LmNvbmZpZwppbXBvcnQgc2FsdC5ydW5uZXIKaW1wb3J0IHNhbHQuY2xpZW50CgpkZWYgZ2V0R3JhaW5zKCBkYywgZW52LCBfX29wdHNfXyApOgoKICAgIHJ1bm5lciA9IHNhbHQucnVubmVyLlJ1bm5lcihfX29wdHNfXykKCiAgICBzdGRvdXRfYmFrID0gc3lzLnN0ZG91dAogICAgd2l0aCBvcGVuKG9zLmRldm51bGwsICd3YicpIGFzIGY6CiAgICAgICAgc3lzLnN0ZG91dCA9IGYKICAgICAgICBpdGVtcyA9IHJ1bm5lci5jbWQoImNhY2hlLmdyYWlucyIsICIiKQogICAgc3lzLnN0ZG91dCA9IHN0ZG91dF9iYWsKCiAgICBhcnIgPSBbXQoKICAgIGlmIGxlbihkYykgPiAwIGFuZCBsZW4oZW52KSA+IDA6CiAgICAgICAgaT0wCiAgICAgICAgZm9yIGtleSBpbiBpdGVtczoKICAgICAgICAgICAgI21hdGNoT2JqID0gcmUubWF0Y2goIGVudiwgaXRlbXNba2V5XVsnZW52J10sIHJlLkkpCiAgICAgICAgICAgICNpZiBpdGVtc1trZXldWydkYyddID09IGRjIGFuZCBtYXRjaE9iajoKICAgICAgICAgICAgI3ByaW50IGl0ZW1zW2tleV1bJ2RjJ10sICBpdGVtc1trZXldWydlbnYnXQogICAgICAgICAgICBpdGVtZGMgPSBpdGVtc1trZXldWydkYyddLmxvd2VyKCkKICAgICAgICAgICAgaXRlbWVudiA9IGl0ZW1zW2tleV1bJ2VudiddLmxvd2VyKCkKICAgICAgICAgICAgZGMgPSBkYy5sb3dlcigpCiAgICAgICAgICAgIGVudiA9IGVudi5sb3dlcigpCiAgICAgICAgICAgIGlmIGl0ZW1kYyA9PSBkYyBhbmQgaXRlbWVudiAgPT0gZW52OgogICAgICAgICAgICAgICAgYXJyLmFwcGVuZCggaXRlbXNba2V5XVsnaWQnXSApCiAgICAgICAgICAgICAgICArK2kKICAgIGVsaWYgbGVuKGRjKSA+IDA6CiAgICAgICAgZm9yIGtleSBpbiBpdGVtczoKICAgICAgICAgICAgaXRlbWRjID0gaXRlbXNba2V5XVsnZGMnXS5sb3dlcigpCiAgICAgICAgICAgIGRjID0gZGMubG93ZXIoKQogICAgICAgICAgICBpZiBpdGVtZGMgPT0gZGM6CiAgICAgICAgICAgICAgICBhcnIuYXBwZW5kKCBpdGVtc1trZXldWydpZCddICkKICAgIGVsc2U6CiAgICAgICAgZm9yIGtleSBpbiBpdGVtczoKICAgICAgICAgICAgYXJyLmFwcGVuZCggaXRlbXNba2V5XVsnaWQnXSApCgogICAgcmV0dXJuIHNvcnRlZChhcnIpLCBpdGVtcwoKZGVmIGdldFNlcnZlcnMoIGRjLCBlbnYsIF9fb3B0c19fICk6CgogICAgbG9jYWwgPSBzYWx0LmNsaWVudC5Mb2NhbENsaWVudCgpCgogICAgaWYgbGVuKGRjKSA+IDAgYW5kIGxlbihlbnYpID4gMDoKICAgICAgICBtYXRjaCA9ICdHQGRjOicgKyBkYyArICcgYW5kIEdAZW52OicgKyBlbnYKICAgIGVsaWYgbGVuKGRjKSA+IDA6CiAgICAgICAgbWF0Y2ggPSAnR0BkYzonICsgZGMKICAgIGVsc2U6CiAgICAgICAgbWF0Y2ggPSAnKicKCiAgICAjIFJlZGlyZWN0IG91dHB1dCB0byAvZGV2L251bGwsIHRvIGVhdCBlcnJvciBvdXRwdXQuCiAgICBzdGRvdXRfYmFrID0gc3lzLnN0ZG91dAogICAgd2l0aCBvcGVuKG9zLmRldm51bGwsICd3YicpIGFzIGY6CiAgICAgICAgc3lzLnN0ZG91dCA9IGYKICAgICAgICAjaXRlbXMgPSBsb2NhbC5jbWQoIG1hdGNoLCAndGVzdC5waW5nJywgW10sIGV4cHJfZm9ybT0nY29tcG91bmQnICkKICAgICAgICBpdGVtcyA9IGxvY2FsLmNtZCggbWF0Y2gsICdncmFpbnMuaXRlbScsIFsiZW52IiwidmVyc2lvbiJdLCBleHByX2Zvcm09J2NvbXBvdW5kJyApCiAgICBzeXMuc3Rkb3V0ID0gc3Rkb3V0X2JhawoKICAgIGFyciA9IFtdCiAgICBpID0gMAoKICAgICMgRmxhdHRlbiB0aGUgb3V0cHV0LiBKdXN0IHdhbnQgWydzcnYxJywnc3J2MicuLi5dCiAgICBmb3Iga2V5IGluIGl0ZW1zOgogICAgICAgIGFyci5hcHBlbmQoIGtleSApCiAgICAgICAgKytpCgogICAgcmV0dXJuIHNvcnRlZChhcnIpLGl0ZW1zIAoKaWYgX19uYW1lX18gPT0gJ19fbWFpbl9fJzoKCiAgICBfX29wdHNfXyA9IHNhbHQuY29uZmlnLm1hc3Rlcl9jb25maWcoCiAgICAgICAgICAgIG9zLmVudmlyb24uZ2V0KCdTQUxUX01BU1RFUl9DT05GSUcnLCAnL2V0Yy9zYWx0L21hc3RlcicpKQoKICAgICMgR2V0IGNvbW1hbmQgbGluZSBhcmdzCgogICAgc2NyaXB0ID0gc3lzLmFyZ3YucG9wKDApCgogICAgZGMgPSAiIgogICAgaWYgbGVuKHN5cy5hcmd2KT4wOgogICAgICAgIGRjID0gc3lzLmFyZ3YucG9wKDApCgogICAgZW52ID0gIiIKICAgIGlmIGxlbihzeXMuYXJndik+MDoKICAgICAgICBlbnYgPSBzeXMuYXJndi5wb3AoMCkKCiAgICAjIEdldCBncmFpbnMgY2FjaGUgaW50byBmbGF0IGFycmF5CgogICAgZ3JhaW5zX2NhY2hlLGNhY2hlX2l0ZW1zID0gZ2V0R3JhaW5zKCBkYywgZW52LCBfX29wdHNfXyApCgogICAgIyBQaW5nIHNlcnZlcnMgdXNpbmcgU2FsdAoKICAgIHNlcnZlcnMsc2VydmVyX2l0ZW1zID0gZ2V0U2VydmVycyggZGMsIGVudiwgX19vcHRzX18gKQoKICAgICMgV2hhdCdzIHRoZSBkaWZmZXJlbmNlCgogICAgc2VydmVyc19zZXQgPSBzZXQoIHNlcnZlcnMgKQogICAgZ3JhaW5zX2NhY2hlX3NldCA9IHNldCggZ3JhaW5zX2NhY2hlICkKICAgIGRpZmZlcmVuY2UgPSBncmFpbnNfY2FjaGVfc2V0IC0gc2VydmVyc19zZXQKCiAgICAjIFNlbmQgYmFjayBhbiBhcnJheSBvZiBvYmplY3RzIGFsbCBzZXQgdXAgbmljZWx5IGZvciBkZXBsb3ltYW4KCiAgICByZXQgPSBbXQogICAgaT0wCiAgICBmb3IgYSBpbiBkaWZmZXJlbmNlOgogICAgICAgIG9iaj17fQogICAgICAgIG9ialsnTmFtZSddPSBhCiAgICAgICAgb2JqWydTZWxlY3RlZCddID0gRmFsc2UKICAgICAgICBvYmpbJ1Jlc3BvbmRlZCddID0gRmFsc2UKICAgICAgICBpZiAndmVyc2lvbicgaW4gY2FjaGVfaXRlbXNbYV06CiAgICAgICAgICAgIG9ialsnVmVyc2lvbiddID0gY2FjaGVfaXRlbXNbYV1bInZlcnNpb24iXQogICAgICAgIGlmICdlbnYnIGluIGNhY2hlX2l0ZW1zW2FdOgogICAgICAgICAgICBvYmpbJ0VudiddID0gY2FjaGVfaXRlbXNbYV1bImVudiJdCiAgICAgICAgcmV0LmFwcGVuZCggb2JqICkKICAgIGZvciBhIGluIHNlcnZlcnM6CiAgICAgICAgb2JqPXt9CiAgICAgICAgb2JqWydOYW1lJ109IGEKICAgICAgICBvYmpbJ1NlbGVjdGVkJ10gPSBGYWxzZQogICAgICAgIG9ialsnUmVzcG9uZGVkJ10gPSBUcnVlCiAgICAgICAgaWYgJ3ZlcnNpb24nIGluIHNlcnZlcl9pdGVtc1thXToKICAgICAgICAgICAgb2JqWydWZXJzaW9uJ10gPSBzZXJ2ZXJfaXRlbXNbYV1bInZlcnNpb24iXQogICAgICAgIGlmICdlbnYnIGluIHNlcnZlcl9pdGVtc1thXToKICAgICAgICAgICAgb2JqWydFbnYnXSA9IHNlcnZlcl9pdGVtc1thXVsiZW52Il0KICAgICAgICByZXQuYXBwZW5kKCBvYmogKQoKICAgIHByaW50IGpzb24uZHVtcHMoIHJldCApCgo="
# }' $proto://$ipport/api/admin/$guid/scripts
# curl -k -d '{
#     "Desc": "Set grains. Arg1 - salt_id, Arg2 - grain,value .. ArgN - grain,value",
#     "Name": "salt-set-grains.sh",
#     "Source": "IyEvYmluL2Jhc2gKCmV4cG9ydCBQQVRIPS9iaW46L3Vzci9iaW4KCmlmICEgd2hpY2ggc2FsdCA+JiAvZGV2L251bGw7IHRoZW4KICAgIGVjaG8gJ3siRXJyb3IiOiJTYWx0IGJpbmFyeSBub3QgZm91bmQuIn0nCiAgICBleGl0IDEKZmkKCltbIC16ICQxIF1dICYmIHsKICBlY2hvICd7IkVycm9yIjoiRmlyc3QgYXJndW1lbnQsIHNhbHRfaWQsIHdhcyBub3Qgc2V0In0nCiAgZXhpdCAxCn0KCnNhbHRfaWQ9JDEKc2hpZnQKCltbIC16ICQxIF1dICYmIHsKICBlY2hvICd7IkVycm9yIjoiTm8gZ3JhaW5zIHdlcmUgc3BlY2lmaWVkIn0nCiAgZXhpdCAxCn0KCmZvciBpIGluICIkQCI7IGRvCiAgY21kPSJzYWx0IC0tb3V0cHV0PWpzb24gJHNhbHRfaWQgZ3JhaW5zLnNldHZhbCAke2klJSwqfSAke2kjIyosfSIKICBvdXRwdXQ9YGVudiAkY21kYAoKICBbWyAteiAkb3V0cHV0IF1dICYmIHsKICAgIGVjaG8gIntcIkVycm9yXCIsXCJObyBvdXRwdXQgZnJvbSBjb21tYW5kLiAkY21kXCIiCiAgICBleGl0IDEKICB9CgogIG51bWxpbmVzPWBlY2hvICIkb3V0cHV0IiB8IHdjIC1sYAogIFtbICRudW1saW5lcyAtbmUgMSBdXSAmJiB7CiAgICBlY2hvICJ7XCJFcnJvclwiLFwiRXhwZWN0ZWQgMSBsaW5lIG9mIG91dHB1dC4gR290ICRudW1saW5lcy5cIn0iCiAgICBlY2hvICJ7XCJPdXRwdXRcIixcIiRvdXRwdXRcIn0iIHwgc2VkIC1uICcxe2g7bn07IHtIfSA7ICQge2c7cy9cbi9cXG4vZztwfScKICAgIGV4aXQgMgogIH0KCiAgZWNobyAiJG91dHB1dCIKZG9uZQoKZXhpdCAwCg=="
# }' $proto://$ipport/api/admin/$guid/scripts
# curl -k -d '{
#     "Desc": "Runs a state.highstate on server(s). Arg1 .. ArgN - salt IDs.",
#     "Name": "salt-highstate.py",
#     "Source": "IyEvdXNyL2Jpbi9lbnYgcHl0aG9uCiMKIyBEZXBsb3ltZW50IE1hbmFnZXIgLSBhIFJFU1QgaW50ZXJmYWNlIGFuZCBHVUkgZm9yIGRlcGxveWluZyBzb2Z0d2FyZQojIENvcHlyaWdodCAoQykgMjAxNCAgTWFyayBDbGFya3NvbgojCiMgVGhpcyBwcm9ncmFtIGlzIGZyZWUgc29mdHdhcmU6IHlvdSBjYW4gcmVkaXN0cmlidXRlIGl0IGFuZC9vciBtb2RpZnkKIyBpdCB1bmRlciB0aGUgdGVybXMgb2YgdGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGFzIHB1Ymxpc2hlZCBieQojIHRoZSBGcmVlIFNvZnR3YXJlIEZvdW5kYXRpb24sIGVpdGhlciB2ZXJzaW9uIDMgb2YgdGhlIExpY2Vuc2UsIG9yCiMgKGF0IHlvdXIgb3B0aW9uKSBhbnkgbGF0ZXIgdmVyc2lvbi4KIwojIFRoaXMgcHJvZ3JhbSBpcyBkaXN0cmlidXRlZCBpbiB0aGUgaG9wZSB0aGF0IGl0IHdpbGwgYmUgdXNlZnVsLAojIGJ1dCBXSVRIT1VUIEFOWSBXQVJSQU5UWTsgd2l0aG91dCBldmVuIHRoZSBpbXBsaWVkIHdhcnJhbnR5IG9mCiMgTUVSQ0hBTlRBQklMSVRZIG9yIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFLiAgU2VlIHRoZQojIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGZvciBtb3JlIGRldGFpbHMuCiMKIyBZb3Ugc2hvdWxkIGhhdmUgcmVjZWl2ZWQgYSBjb3B5IG9mIHRoZSBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZQojIGFsb25nIHdpdGggdGhpcyBwcm9ncmFtLiAgSWYgbm90LCBzZWUgPGh0dHA6Ly93d3cuZ251Lm9yZy9saWNlbnNlcy8+LgoKaW1wb3J0IGpzb24KaW1wb3J0IG9zCmltcG9ydCBzeXMKaW1wb3J0IGdldG9wdAppbXBvcnQgcmUKCmltcG9ydCBzYWx0LmNvbmZpZwppbXBvcnQgc2FsdC5ydW5uZXIKaW1wb3J0IHNhbHQuY2xpZW50CgpkZWYgZG9IaWdoU3RhdGUoIHNlcnZlcnMsIF9fb3B0c19fICk6CgogICAgbG9jYWwgPSBzYWx0LmNsaWVudC5Mb2NhbENsaWVudCgpCgogICAgam9iaWRzID0gW10KCiAgICAjIFJlZGlyZWN0IG91dHB1dCB0byAvZGV2L251bGwsIHRvIGVhdCBlcnJvciBvdXRwdXQuCiAgICBzdGRvdXRfYmFrID0gc3lzLnN0ZG91dAogICAgd2l0aCBvcGVuKG9zLmRldm51bGwsICd3YicpIGFzIGY6CiAgICAgICAgc3lzLnN0ZG91dCA9IGYKICAgICAgICAjaXRlbXMgPSBsb2NhbC5jbWRfaXRlcl9ub19ibG9jayggc2VydmVycywgJ3Rlc3QucGluZycsIFtdLCBleHByX2Zvcm09J2xpc3QnICkKICAgICAgICBmb3IgaSBpbiBzZXJ2ZXJzOgogICAgICAgICAgICBqb2JpZHMuYXBwZW5kKCB7ICdTYWx0SWQnOmksICdKb2JJZCc6bG9jYWwuY21kX2FzeW5jKGksJ3N0YXRlLmhpZ2hzdGF0ZScsW10pfSApCiAgICBzeXMuc3Rkb3V0ID0gc3Rkb3V0X2JhawoKICAgIHJldHVybiBqb2JpZHMgCgppZiBfX25hbWVfXyA9PSAnX19tYWluX18nOgoKICAgIF9fb3B0c19fID0gc2FsdC5jb25maWcubWFzdGVyX2NvbmZpZygKICAgICAgICAgICAgb3MuZW52aXJvbi5nZXQoJ1NBTFRfTUFTVEVSX0NPTkZJRycsICcvZXRjL3NhbHQvbWFzdGVyJykpCgogICAgIyBHZXQgY29tbWFuZCBsaW5lIGFyZ3MKCiAgICBzY3JpcHQgPSBzeXMuYXJndi5wb3AoMCkKCiAgICBzZXJ2ZXJzID0gW10gCiAgICBpZiBsZW4oc3lzLmFyZ3YpPjA6CiAgICAgICAgd2hpbGUgbGVuKHN5cy5hcmd2KT4wOgogICAgICAgICAgICBzZXJ2ZXJzLmFwcGVuZCggc3lzLmFyZ3YucG9wKDApICkKICAgIGVsc2U6CiAgICAgICAgcHJpbnQgIk5vIHNlcnZlcnMgd2VyZSBzcGVjaWZpZWQiCiAgICAgICAgc3lzLmV4aXQoMSkKCiAgICAjIFB1dCB0aGUgc2VydmVycyBpbnRvIHN0YXRlLmhpZ2hzdGF0ZQoKICAgIGpvYmlkcyA9IGRvSGlnaFN0YXRlKCBzZXJ2ZXJzLCBfX29wdHNfXyApCgogICAgIyBTZW5kIGJhY2sgYW4gYXJyYXkgb2Ygb2JqZWN0cyBhbGwgc2V0IHVwIG5pY2VseSBmb3IgZGVwbG95bWFuCgogICAgcHJpbnQganNvbi5kdW1wcyggam9iaWRzICkKCg=="
# }' $proto://$ipport/api/admin/$guid/scripts

#' ----

curl -k -d '{
    "Name":"yum",
    "Desc":"Yum main sidebar item",
    "HasView":2,
    "Parent":""
}' $proto://$ipport/api/admin/$guid/plugins

curl -k -d '{
    "Name":"yumupdate",
    "Desc":"Update the YUM repo from Jenkins",
    "HasView":1,
    "Parent":"yum"
}' $proto://$ipport/api/admin/$guid/plugins

curl -k -d '{
    "Name":"systemlog",
    "Desc":"View all system activity",
    "HasView":1,
    "Parent":""
}' $proto://$ipport/api/admin/$guid/plugins | tee /tmp/fshodifha
# Grab the id of the last insert
id=`grep Id /tmp/fshodifha | grep -Eo "[0-9]+"`
# Add the AJS controller file
curl -k -d '{
    "Name":"systemlog.js",
    "Desc":"Controller for the system log",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"systemlog/js/controllers/systemlog.js"
}' $proto://$ipport/api/admin/$guid/files

curl -k -d '{
    "Name":"saltkeymanager",
    "Desc":"Salt key management plugin",
    "HasView":1,
    "Parent":"salt"
}' $proto://$ipport/api/admin/$guid/plugins | tee /tmp/fshodifha
# Grab the id of the last insert
id=`grep Id /tmp/fshodifha | grep -Eo "[0-9]+"`
# Add the AJS controller file
curl -k -d '{
    "Name":"saltkeymanager.js",
    "Desc":"Controller for Salt key manager",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"saltkeymanager/js/controllers/saltkeymanager.js"
}' $proto://$ipport/api/admin/$guid/files

echo "GUID=$guid"
