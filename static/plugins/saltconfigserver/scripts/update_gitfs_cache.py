#!/bin/env python
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
'''
Fills the gitfs fileserver cache.
Expects one argument, the environment name, e.g. dev_0.1.2
'''

# Import python libs
import os
import sys

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
        load = {'saltenv':env}
        file = salt.fileserver.gitfs.file_list(load)
        for i in file:
            result.append( salt.fileserver.gitfs.find_file(i,env) )
    except KeyboardInterrupt:
        os.kill(pid, 15)

    return result

if __name__ == '__main__':
    sys.argv.pop(0)
    env = ""
    if len(sys.argv)>0:
        env = sys.argv.pop(0)
    else:
        print '{"Error":"Arg1, the environment, is empty."}'
        sys.exit(1)

    print salt_update()
