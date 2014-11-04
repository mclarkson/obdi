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
    "Name":"jobs",
    "Desc":"Job management",
    "HasView":1,
    "Parent":""
}' $proto://$ipport/api/admin/$guid/plugins | tee $t

# Grab the id of the last insert
id=`grep Id $t | grep -Eo "[0-9]+"`

# Delete the temporary file and delete the trap
rm -f -- "$t"
trap - EXIT

#
# Add the AJS controller files
#
# These need to be loaded when the application starts
#

curl -k -d '{
    "Name":"jobs.js",
    "Desc":"Manages the job queue",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"jobs/js/controllers/jobs.js"
}' $proto://$ipport/api/admin/$guid/files
curl -k -d '{
    "Name":"outputlines.js",
    "Desc":"Shows job log",
    "Type":1,
    "PluginId":'"$id"',
    "Url":"jobs/js/controllers/outputlines.js"
}' $proto://$ipport/api/admin/$guid/files

