#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # byte(0x01) -> text format
    SENT_DATA=$(echo -ne "\x01hello, world!" | base64);
else
    echo "Sending file..." && echo # byte(0x02) -> file format
    SENT_DATA=$(echo -ne "\x02example.txt\x02hello, world!" | base64);
fi

PUSH_FORMAT='{
    "receiver":"Alice",
    "req_data":{
        "method":"POST",
        "host":"hidden-lake-messenger",
        "path":"/push",
        "body":"'${SENT_DATA}'"
    }
}';

curl -i -X PUT -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo 
