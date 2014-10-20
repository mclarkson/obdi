#!/bin/bash

proto="https"
opts="-k -s" # don't check ssl cert, silent
ipport="127.0.0.1:4443"

curl -k -d '{ 
    "ScriptSource":"IyEvYmluL2Jhc2gKClBBVEg9JFBBVEg6L2JpbjovdXNyL2JpbgoKZWNobyAiSGkiCgpzbGVlcCAxCgplY2hvICRACgpzbGVlcCAxCgplbnYKCnNsZWVwIDEKCmVjaG8gIkJ5ZSIKCmV4aXQgMAoK",
    "EnvVars":"A=1 B=\"Hi There\" C=3",
    "NotifURL":"https://127.0.0.1",
    "JobID":123,
    "Key":"lOcAlH0St",
    "Args":"-a \"hi there\" -x 1"
}' $proto://$ipport/api/jobs

