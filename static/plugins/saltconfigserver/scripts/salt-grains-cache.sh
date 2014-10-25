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

salt_id=$1

# The server is offline, so don't:
# #output=`salt --output=json "$salt_id" grains.items`
output=

[[ -z $output ]] && {

output=`python < <(cat <<EnD
#!/usr/bin/env python
import json
import os
import sys
import getopt

import salt.config
import salt.runner

if __name__ == '__main__':
    __opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/minion'))
    runner = salt.runner.Runner(__opts__)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        items = runner.cmd("cache.grains", ["$salt_id"])
    sys.stdout = stdout_bak

    print json.dumps(items)
EnD
)`
}

echo $output
