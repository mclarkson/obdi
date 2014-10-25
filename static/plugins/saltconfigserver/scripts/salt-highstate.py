#!/usr/bin/env python
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

import json
import os
import sys
import getopt
import re

import salt.config
import salt.runner
import salt.client

def doHighState( servers, __opts__ ):

    local = salt.client.LocalClient()

    jobids = []

    # Redirect output to /dev/null, to eat error output.
    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        #items = local.cmd_iter_no_block( servers, 'test.ping', [], expr_form='list' )
        for i in servers:
            jobids.append( { 'SaltId':i, 'JobId':local.cmd_async(i,'state.highstate',[])} )
    sys.stdout = stdout_bak

    return jobids 

if __name__ == '__main__':

    __opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/master'))

    # Get command line args

    script = sys.argv.pop(0)

    servers = [] 
    if len(sys.argv)>0:
        while len(sys.argv)>0:
            servers.append( sys.argv.pop(0) )
    else:
        print "No servers were specified"
        sys.exit(1)

    # Put the servers into state.highstate

    jobids = doHighState( servers, __opts__ )

    # Send back an array of objects all set up nicely for obdi

    print json.dumps( jobids )

