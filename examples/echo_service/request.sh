#!/bin/bash

str2hex() {
    local str=${1:-""}
    local fmt="%02X"
    local chr
    local -i i
    for i in `seq 0 $((${#str}-1))`; do
        chr=${str:i:1}
        printf "${fmt}" "'${chr}"
    done
}

JSON_DATA='{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
            "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
}';

PUSH_FORMAT='{
        "receiver":"Bob",
        "hex_data":"'$(str2hex "$JSON_DATA")'"
}';

d="$(date +%s)";
curl -i -X POST -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
