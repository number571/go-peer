#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # encrypted(byte(0x01) -> text format)
    SENT_DATA=$(echo -ne "\xf6\xa5\xc2\xc8\x22\x22\xd2\x50\x37\xf6\x51\xa5\xd5\x33\xd2\xfa\x1b\xc7\x5a\x2d\xf0\xd6\xda\x6d\x27\x6f\xed\x36\x12\xdf\x7d\xc4\xa5\x7f\xdf\xdc\x1f\x57\xb8\xc1\x47\xba\x3e\xc2\x24\x72\xe0\x3f\xf9\x4c\x3f\x12\xa6\x1b\xc7\x10\x6e\xe1\x5f\xe7\x3b\x18" | base64);
else
    echo "Sending file..." && echo # encrypted(byte(0x02) -> file format)
    SENT_DATA=$(echo -ne "\xbf\x13\xc0\x13\x6f\xb0\xc2\x86\xfc\x84\x1e\xc9\x7c\xa3\x6b\xe7\xf0\xc3\xfd\x15\x18\x0f\x17\x9a\xcb\x0a\x9c\x72\x4a\x8f\x3d\x0c\xe5\xc3\x1c\xd7\xdb\xad\x41\x0d\xde\xef\xa6\xff\x5e\x01\x66\xc7\x6b\x14\x21\x1d\xf4\x98\x26\x96\x0a\xf6\x47\xc5\xd5\xd2\xa3\x2c\x50\x68\xa0\x59\xea\x14\x9d\xea\x19\x11" | base64);
fi

REQUEST_ID=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 44)
JSON_DATA='{
        "method":"POST",
        "host":"hidden-lake-messenger",
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
