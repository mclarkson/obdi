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
    "Name":"saltupdategit",
    "Desc":"Update the Salt Git repository",
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
    "Name":"saltupdategit.js",
    "Desc":"Controller for Salt Update Git Repo plugin",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"saltupdategit/js/controllers/saltupdategit.js"
}' $proto://$ipport/api/admin/$guid/files

#
# Add the scripts, removing comment lines (#) and empty lines
#

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/saltupdategit_set_version.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=saltupdategit_set_version.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
			"Desc": "Returns a list of versions (branches). Arg1 - The branch to search on.",
		"Name": "saltupdategit_set_version.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/saltupdategit_get_version.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=saltupdategit_get_version.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
			"Desc": "Returns a list of versions (branches). Arg1 - The branch to search on.",
		"Name": "saltupdategit_get_version.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# --

source=`sed '1n;/^\s*#/d;/^$/d;' scripts/saltupdategit_increment_version.sh | base64 -w 0`

curl -k $proto://$ipport/api/admin/$guid/scripts?name=saltupdategit_increment_version.sh | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

if [[ -z $id ]]; then
	curl -k -d '{
		"Desc": "Increment the Salt Git version creating a new branch from HEAD. Arg1 - The branch name, Arg2 - The position to increment, 1, 2 or 3.",
		"Name": "saltupdategit_increment_version.sh",
		"Source": "'"$source"'"
	}' $proto://$ipport/api/admin/$guid/scripts
else
	curl -k -X PUT -d '{ "Source": "'"$source"'" }' \
	$proto://$ipport/api/admin/$guid/scripts/$id
fi

# Delete the temporary file and delete the trap
rm -f -- "$t"
trap - EXIT

