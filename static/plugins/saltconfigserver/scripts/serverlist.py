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

def getGrains( dc, env, __opts__ ):

    runner = salt.runner.Runner(__opts__)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        items = runner.cmd("cache.grains", "")
    sys.stdout = stdout_bak

    arr = []

    if len(dc) > 0 and len(env) > 0:
        i=0
        for key in items:
            #matchObj = re.match( env, items[key]['env'], re.I)
            #if items[key]['dc'] == dc and matchObj:
            #print items[key]['dc'],  items[key]['env']
            if 'dc' in items[key]:
                itemdc = items[key]['dc'].lower()
            else:
                continue
            if 'env' in items[key]:
                itemenv = items[key]['env'].lower()
            else:
                continue
            dc = dc.lower()
            env = env.lower()
            if itemdc == dc and itemenv  == env:
                arr.append( items[key]['id'] )
                ++i
    elif len(dc) > 0:
        for key in items:
            if 'dc' in items[key]:
                itemdc = items[key]['dc'].lower()
            else:
                continue
            dc = dc.lower()
            if itemdc == dc:
                arr.append( items[key]['id'] )
    else:
        for key in items:
            arr.append( items[key]['id'] )

    return sorted(arr), items

def getServers( dc, env, __opts__ ):

    local = salt.client.LocalClient()

    if len(dc) > 0 and len(env) > 0:
        match = 'G@dc:' + dc + ' and G@env:' + env
    elif len(dc) > 0:
        match = 'G@dc:' + dc
    else:
        match = '*'

    # Redirect output to /dev/null, to eat error output.
    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        #items = local.cmd( match, 'test.ping', [], expr_form='compound' )
        items = local.cmd( match, 'grains.item', ["env","version"], expr_form='compound',timeout=10 )
    sys.stdout = stdout_bak

    arr = []
    i = 0

    # Flatten the output. Just want ['srv1','srv2'...]
    for key in items:
        arr.append( key )
        ++i

    return sorted(arr),items 

if __name__ == '__main__':

    __opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/master'))

    # Get command line args

    script = sys.argv.pop(0)

    dc = ""
    if len(sys.argv)>0:
        dc = sys.argv.pop(0)

    env = ""
    if len(sys.argv)>0:
        env = sys.argv.pop(0)

    # Get grains cache into flat array

    grains_cache,cache_items = getGrains( dc, env, __opts__ )

    # Ping servers using Salt

    servers,server_items = getServers( dc, env, __opts__ )

    # What's the difference

    servers_set = set( servers )
    grains_cache_set = set( grains_cache )
    difference = grains_cache_set - servers_set

    # Send back an array of objects all set up nicely for obdi

    ret = []
    i=0
    for a in difference:
        obj={}
        obj['Name']= a
        obj['Selected'] = False
        obj['Responded'] = False
        if 'version' in cache_items[a]:
            obj['Version'] = cache_items[a]["version"]
        if 'env' in cache_items[a]:
            obj['Env'] = cache_items[a]["env"]
        ret.append( obj )
    for a in servers:
        obj={}
        obj['Name']= a
        obj['Selected'] = False
        obj['Responded'] = True
        if 'version' in server_items[a]:
            obj['Version'] = server_items[a]["version"]
        if 'env' in server_items[a]:
            obj['Env'] = server_items[a]["env"]
        ret.append( obj )

    print json.dumps( ret )

