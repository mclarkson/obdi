#!/bin/bash

[[ -z $1 ]] && {
    echo '{ "Error":"Must specify the name as Arg1." }'
    exit 1
}

salt-key -y --output=json -r $1
