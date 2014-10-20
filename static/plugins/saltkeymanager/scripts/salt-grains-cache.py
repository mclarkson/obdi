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

    script = sys.argv.pop(0)
    dc = ""
    env = ""

    if len(sys.argv)>0:
        dc = sys.argv.pop(0)
    if len(sys.argv)>0:
        env = sys.argv.pop(0)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        items = runner.cmd("cache.grains", "")
    sys.stdout = stdout_bak

    ret = {}

    if len(dc) > 0 and len(env) > 0:
        for key in items:
            if items[key]['dc'] == dc and items[key]['env'] == env:
                ret[key] =  items[key]
    elif len(dc) > 0:
        for key in items:
            if items[key]['dc'] == dc:
                ret[key] =  items[key]
    else:
        for key in items:
            ret[key] =  items[key]

    print json.dumps(ret)

