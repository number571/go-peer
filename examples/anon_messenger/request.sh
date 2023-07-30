#!/bin/bash

## Works only if users are logged in to the account!
## node2[localhost:7070] -> node1[localhost:8080]

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

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # byte(0x01) -> text format
    SENT_DATA=$(echo -ne "\x01hello, world!" | base64);
else
    echo "Sending file..." && echo # byte(0x02) -> file format
    SENT_DATA=$(echo -ne "\x02example.txt\x02hello, world!" | base64);
fi

JSON_DATA='{
        "method":"POST",
        "host":"go-peer/hidden-lake-messenger",
        "path":"/push",
        "head":{
            "Accept": "application/json"
        },
        "body":"'${SENT_DATA}'"
}';

PUSH_FORMAT='{
        "receiver":"Alice",
        "hex_data":"'$(str2hex "$JSON_DATA")'"
}';

curl -i -X PUT -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo 
