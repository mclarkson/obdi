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

export PATH=/bin:/usr/bin

if ! which salt >& /dev/null; then
    echo '{"Error":"Salt binary not found."}'
    exit 1
fi

[[ -z $1 ]] && {
    echo '{ "Error":"Arg1 empty. Expected environment name" }'
    exit 1
}

ENV=$1

# Does the tag/branch exist?

# Get the gitfs remote from the salt master file
gitrepo=`grep -A1 gitfs_remotes: /etc/salt/master | sed -n '${s/^\s*-\s*//p}'`

# Get a list of all the branches/tags
branchtags=`git ls-remote $gitrepo | sed -n '/HEAD/n;s#^.*refs/heads/##p'`

if ! echo "$branchtags" | grep -qs "^$ENV\$" 2>/dev/null; then
    echo '{ "Error":"Branch or tag \"'"$ENV"'\" does not exist." }'
    exit 1
fi

# Clear the cache first in case sls files were deleted at source. The
# deleted files will hang around forever otherwise.

[[ -d /var/cache/salt/master/gitfs/refs/$ENV ]] && {
    rm -rf /var/cache/salt/master/gitfs/refs/$ENV
}

# Fill the salt git cache

output=`python < <(cat <<EnD
# Import python libs
import os
import sys
import fnmatch

# Import salt libs
import salt
import salt.cli
import salt.config
import salt.fileserver.gitfs

def salt_update():

    salt.fileserver.gitfs.__opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/master'))

    pid = os.getpid()
    result = []
    try:
        salt.fileserver.gitfs.update()
        load = {'saltenv':'$ENV'}
        file = salt.fileserver.gitfs.file_list(load)
        for i in file:
            if fnmatch.fnmatch(i, "*.sls"):
                result.append( salt.fileserver.gitfs.find_file(i,'$ENV') )
    except KeyboardInterrupt:
        os.kill(pid, 15)

    return result

print salt_update()
EnD
)`

# Does the env exist now?

[[ ! -d /var/cache/salt/master/gitfs/refs/$ENV ]] && {
    echo '{ "Error":"Unable to fill git cache", "Output":"'"$output"'" }'
    exit 1
}

# Get sls details and output as JSON

cd /var/cache/salt/master/gitfs/refs/; cd ${ENV}
declare -i numstates i
numstates=`find . -type f -name "*.sls" | wc -l`
echo -n "["
i=0
find . -type f -name "*.sls" | \
    while read line; do \
        class="";
        formula="";
        statefile="";
        [[ "$line" == *"/init.sls" ]] && {
            formula=`echo $line | sed 's#^./\(.*\)/.*#\1#'`;
        } || {
            class=`echo $line | sed 's#^./\(.*\)/\(.*\)\.sls#\1.\2#'`;
                   formula=${class%.*};statefile=${class#*.};
        };
        desc=`head "$line" | sed -n 's/^# INFO: *\(.*\)/\1/p' | \
              tr -d '\0-\37'`;
        echo -n '{"Desc":"'"$desc"'","FormulaName":"'"${formula}";
        echo -n '","StateFileName":"'"${statefile}"'"}';
        i+=1;
        [[ $i -ne $numstates ]] && echo -n ",";
    done
echo "]"

exit 0
