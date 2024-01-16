#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # encrypted(byte(0x01) -> text format)
    REQUEST_ID="Hdvl9TuoTfp3L-0HbsFb2J5tcDnuN0iHgHzbtCrKLRG1"
    SENT_DATA=$(echo -ne "\xcf\x5e\x8f\x27\xd9\x42\xc8\x04\xb4\xdf\xb7\xff\xeb\x20\xe4\x4e\x5a\x16\xe0\xe4\xa2\x9a\x2e\x40\x35\x5d\xc4\x70\xc4\xd6\x33\xb5\xb4\x8c\x71\x3f\x81\xf9\x7a\xe0\x4c\x7f\x2a\x13\xf9\x51\xda\x46\xa4\xb3\xc2\xc5\x98\xe7\x32\x17\x67\xc7\x43\xb1\xc5\x4e" | base64);
else
    echo "Sending file..." && echo # encrypted(byte(0x02) -> file format)
    REQUEST_ID="HrgMvCw10XohGs6gaFU7R44N7QDRwVSw4FPS7rwFcotQ"
    SENT_DATA=$(echo -ne "\xb7\xd6\x39\xb4\x3d\x5b\xb1\x54\x8f\xfc\xf5\x35\xde\x92\x76\x69\x2e\x49\xd3\x62\x52\xaa\x31\x49\xe2\x87\xb6\xaf\xc1\xbc\x0f\x4e\xf7\x26\x19\x4e\x79\x1c\x2b\x32\x54\x87\x56\xd9\xf5\xb3\xfb\x0a\xe9\x1c\x4f\x2d\x32\xbe\x41\xa7\x42\x91\xc3\xef\xf0\xaa\x05\x0f\x09\xa2\xa9\x4f\xff\x0c\xd2\x02\xd4\x9e" | base64);
fi

JSON_DATA='{
        "method":"POST",
        "host":"hidden-lake-messenger",
        "path":"/push",
        "head":{
            "Hl-Messenger-Sender-Id": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
            "Hl-Messenger-Request-Id": "'${REQUEST_ID}'",
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
