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

import urllib
import urllib2
import json
import sys
import time
import fcntl

import os
import re

import salt.config
import salt.runner
import salt.client

def getGrains( g, __opts__ ):

    local = salt.client.LocalClient()

    # Redirect output to /dev/null, to eat error output.
    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        #items = local.cmd( match, 'test.ping', [], expr_form='compound' )
        items = local.cmd( g['salt_id'], 'grains.item', ["dc","env","version"] )
    sys.stdout = stdout_bak

    if 'dc' in items[g['salt_id']]:
        dc = items[g['salt_id']]['dc']
    else:
        dc = 'null'
    if 'env' in items[g['salt_id']]:
        env = items[g['salt_id']]['env']
    else:
        env = 'null'
    if 'version' in items[g['salt_id']]:
        version = items[g['salt_id']]['version']
    else:
        version = 'null'

    return dc,env,version

def getGrainsFromCache( g, __opts__ ):

    runner = salt.runner.Runner(__opts__)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        items = runner.cmd("cache.grains", [g['salt_id']])
    sys.stdout = stdout_bak

    dc = ""
    env = ""
    version = ""

    if g['salt_id'] in items:
        if 'dc' in items[g['salt_id']]:
            dc = items[g['salt_id']]['dc']
        else:
            dc = 'null'
        if 'env' in items[g['salt_id']]:
            env = items[g['salt_id']]['env']
        else:
            env = 'null'
        if 'version' in items[g['salt_id']]:
            version = items[g['salt_id']]['version']
        else:
            version = 'null'

    return dc,env,version

def login( g ):
    url = g['master_url'] + "/api/login"

    creds = { "Login":g['enc_user'],"Password":g['enc_pass'] }

    data = json.dumps(creds)

    req = urllib2.Request(url, data)

    try:
        response = urllib2.urlopen( req )
    except urllib2.HTTPError as e:
        a = json.loads( e.read() )
        print "Server said:", e.code, a['Error']
        sys.exit(1)

    html = json.loads( response.read() )

    return html['GUID']

def enc_query( g ):
    url = ( g['master_url'] + "/api/" + g['enc_user'] + "/"+g['guid'] +
          "/saltconfigserver/enc?env=" + g['env'] +
          "&version=" + g['version'] +
          "&dc=" + g['dc'] +
          "&salt_id=" + g['salt_id'] +
          "&yaml=true" )

    req = urllib2.Request(url)

    try:
        response = urllib2.urlopen( req )
    except urllib2.HTTPError as e:
        a = json.loads( e.read() )
        print "Server said:", e.code, a['Error']
        sys.exit(1)

    html = json.loads( response.read() )

    return html['EncData']

def usage():
    print "enc_query.py salt_id"

def main():

    __opts__ = salt.config.master_config(
        os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/master'))

    # Read the worker configuration file
    conf = {}
    with open("/etc/obdi-worker/obdi-worker.conf") as myfile:
        for line in myfile:
            name, var = line.partition("=")[::2]
            conf[name.strip()] = var.strip(" \"\n")

    # Initialise 'g'
    g = {}

    # From obdi_worker.conf
    g['master_url'] = conf['man_urlprefix']
    g['enc_user'] = conf['man_user']
    g['enc_pass'] = conf['man_password']

    # First item is the script name
    dummy = sys.argv.pop(0)

    if len(sys.argv)>0:
        g['salt_id'] = sys.argv.pop(0)
    else:
        usage()
        sys.exit(1)

    dc,env,version = getGrainsFromCache( g, __opts__ )

    g['dc'] = dc
    g['env'] = env
    g['version'] = version

    # Lock access so GUID isn't invalidated by another process
    # logging in at the same time

    done = False
    tries = 0
    lock_file = '/var/tmp/enc_query.py.lock'

    while done == False:
        fp = open(lock_file, 'w')
        try:
            fcntl.lockf(fp, fcntl.LOCK_EX | fcntl.LOCK_NB)
            done = True
        except IOError:
            # another instance is running
            if tries > 100:
                print "Could not gain lock on " + lock_file
                sys.exit(1)
            tries += 1
            time.sleep( 0.1 )

    g['guid'] = login( g )

    resp = enc_query( g )

    # Unlock the lockfile
    fcntl.lockf(fp, fcntl.LOCK_UN)
    fp.close()

    print resp

main()

