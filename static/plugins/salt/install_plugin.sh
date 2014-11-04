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

# salt stub

curl -k -d '{
    "Name":"salt",
    "Desc":"Salt main sidebar item",
    "HasView":2,
    "Parent":""
}' $proto://$ipport/api/admin/$guid/plugins

