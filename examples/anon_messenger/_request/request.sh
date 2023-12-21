#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # encrypted(byte(0x01) -> text format)
    SENT_DATA=$(echo -ne "\xbf\x34\x1e\x0f\xcd\xbd\xe3\x81\x49\x43\xde\x4b\xb3\x21\x72\xc6\x12\x87\xa9\x6b\x6c\x1d\x8e\xf6\xad\xd1\x29\xd5\xb5\x1c\x16\xd9\x53\x04\x69\x06\x32\xb8\x17\xcc\xc4\xbe\x57\x4f\xc6\xa1\x7d\x98\x60\x5c\x7a\xdc\x31\xed\x1f\xa3\x75\xba\x65\x61\xfd\x8f" | base64);
else
    echo "Sending file..." && echo # encrypted(byte(0x02) -> file format)
    SENT_DATA=$(echo -ne "\x5e\xee\xc2\x4e\x3f\x42\x2d\xf6\x38\xdb\x7a\xbe\x95\x49\x78\x5b\x18\x42\xd3\x78\x9f\xbf\xc3\x53\x28\xca\xae\xe4\x7c\x2d\x91\x62\xdb\xe8\x5d\xa2\xc2\xfa\x1b\xf8\x20\xc1\xcb\x71\x3e\xe7\x82\x76\xb8\xb8\xa7\xc4\xeb\xe9\xaf\x76\x84\xba\xdd\xaa\x3a\xdb\xc5\xd2\x32\x92\xd6\x7e\x05\x8f\xce\x30\x28\xc2" | base64);
fi

REQUEST_ID=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 16)
JSON_DATA='{
        "method":"POST",
        "host":"go-peer/hidden-lake-messenger",
        "path":"/push",
        "head":{
            "Hl-Messenger-Sender-Id": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
            "Hl-Service-Request-Id": "'${REQUEST_ID}'",
            "Accept": "application/json"
        },
        "body":"'${SENT_DATA}'"
}';

JSON_DATA=${JSON_DATA//\"/\\\"} # "method" -> \"method\", ...
JSON_DATA=${JSON_DATA//[$'\t\r\n ']} # delete \t \r \n ' ' from string

PUSH_FORMAT='{
        "receiver":"Alice",
        "req_data":"'$JSON_DATA'"
}';

curl -i -X PUT -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo 
