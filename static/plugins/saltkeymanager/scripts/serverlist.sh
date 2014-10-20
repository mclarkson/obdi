#!/bin/bash

#salt --output=json -C "G@dc:$1 and G@env:$2" grains.items

echo -n "[";
salt --output=json -C "G@dc:$1 and G@env:$2" test.ping |
    while read line
    do
        echo $line, | tr -d '\n'
    done | sed '$ {s/,$//}'
echo -n "]"

