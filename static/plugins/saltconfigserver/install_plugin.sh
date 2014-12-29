#!/bin/bash
#
# Obdi - a REST interface and GUI for deploying software
# Copyright (C) 2014  Mark Clarkson
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
# saltconfigserver plugin

#
# Arg1 - the name of the root owned, 0600 permission, file with the
#        json encoded data for initial login, i.e.
#
#            {"Login":"admin","Password":"admin"}
#

# TODO : read pw file
[[ -z $1 ]] && {
    echo "ERROR: Arg1 - name of file containing json credentials data."
    exit 1
}

[[ ! -r "$1" ]] && {
    echo "ERROR: Could not find file '$1'. Aborting"
    exit 2
}

pwfile="$1"

##
## TODO: This plugin depends on salt
##

proto="https"
opts="-k -s" # don't check ssl cert, silent
ipport="127.0.0.1:443"
guid=`curl $opts -f -d @$pwfile \
    $proto://$ipport/api/login | grep -o "[a-z0-9][^\"]*"`

[[ $? -ne 0 ]] && {
    curl $opts -s -d @$pwfile $proto://$ipport/api/login
    echo "Login error"
    exit 1
}

echo "GUID=$guid"

#
# Create a temporary file and a trap to delete it
#

t="/tmp/install_saltconfigserver_$$"
touch $t
[[ $? -ne 0 ]] && {
    echo "Could not create temporary file. Aborting."
    exit 1
}
trap "rm -f -- '$t'" EXIT

#
# Create the plugin entry in obdi, so it can be shown in the sidebar
#

curl -k -d '{
    "Name":"saltconfigserver",
    "Desc":"Salt Configure Server plugin",
    "HasView":1,
    "Parent":"salt"
}' $proto://$ipport/api/admin/$guid/plugins | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

#
# Add the AJS controller files
#
# These need to be loaded when the application starts
#

curl -k -d '{
    "Name":"saltconfigserver.js",
    "Desc":"Controller for Salt New Server",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"saltconfigserver/js/controllers/saltconfigserver.js"
}' $proto://$ipport/api/admin/$guid/files

#
# Add the scripts, removing comment lines (#) and empty lines
#

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/salt-grains-cache.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=salt-grains-cache.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Return the grains for a server from the Salt Cache. Arg1 - Salt ID",
		"Name": "salt-grains-cache.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/salt-grains.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=salt-grains.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Return the grains for a server. Arg1 - Salt ID",
		"Name": "salt-grains.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/serverlist.py | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=serverlist.py | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Returns a list of servers, dead or alive. Arg 1 - dc, Arg 2 - env",
		"Name": "serverlist.py",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/salt-set-grains.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=salt-set-grains.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Set grains. Arg1 - salt_id, Arg2 - grain,value .. ArgN - grain,value",
		"Name": "salt-set-grains.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/salt-highstate.py | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=salt-highstate.py | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Runs a state.highstate on server(s). Arg1 .. ArgN - salt IDs.",
		"Name": "salt-highstate.py",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/statelist.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=statelist.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Returns a list of states and formulas. Arg1 - The branch/tag.",
		"Name": "statelist.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/saltconfigserver_get_version.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=saltconfigserver_get_version.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
			"Desc": "Returns a list of versions (branches). Arg1 - The branch to search on.",
		"Name": "saltconfigserver_get_version.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# Delete the temporary file and delete the trap
rm -f -- "$t"
trap - EXIT

