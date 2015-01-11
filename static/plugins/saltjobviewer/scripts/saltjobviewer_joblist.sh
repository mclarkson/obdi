#!/bin/bash

export PATH=/bin:/usr/bin

# ---------------------------------------------------------------------------
# Salt-run to output json
# ---------------------------------------------------------------------------

salt_run_json() {

    local fun="$1"
    local arg="$2"

    output=`python < <(cat <<EnD
import json
import os
import sys

import salt.config
import salt.runner

if __name__ == '__main__':
    __opts__ = salt.config.master_config(
            os.environ.get('SALT_MASTER_CONFIG', '/etc/salt/minion'))
    runner = salt.runner.Runner(__opts__)

    stdout_bak = sys.stdout
    with open(os.devnull, 'wb') as f:
        sys.stdout = f
        ret = runner.cmd("$fun", ["$arg"])
    sys.stdout = stdout_bak

    print json.dumps(ret)
EnD
    )`

    echo "$output"
}

# ---------------------------------------------------------------------------
# Job check
# ---------------------------------------------------------------------------

# Get minion name (only one is expected and checked for)

result=$(salt_run_json jobs.list_jobs)

echo "$result"

exit 0

