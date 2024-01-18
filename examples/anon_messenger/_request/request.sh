#!/bin/bash

## this script only works once!
## nodes remember the generated request_id

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

randVal=$((RANDOM % 2))
if [ "$randVal" -eq "0" ]; then
    echo "Sending text..." && echo # encrypted(byte(0x01) -> text format)
    REQUEST_ID="23GLKD8ovHnuTv99mMpIxq0O2srtuzUPLwRwq6Kd5VDo"
    SENT_DATA=$(echo -ne "\xfb\xaa\x5a\xb9\x8b\x69\x3e\x58\x16\xdd\x83\x81\x58\xc4\xea\x7f\x62\x25\x1d\x02\xb1\xa5\x4f\xc6\x79\x68\xb9\xbd\x47\x76\xf3\x97\xd5\xd4\xc6\xc4\x15\x1e\x40\x6d\x3b\xb6\x4c\x5c\xa8\xa9\x20\xdc\xec\x52\x78\xd0\xfc\x4a\x74\x99\x72\xca\x4d\x1c\x77\xcf" | base64);
else
    echo "Sending file..." && echo # encrypted(byte(0x02) -> file format)
    REQUEST_ID="HoiAiKaDCPc7V1WKHpCVsDYDLqbtZMramZyDJD-ETMYU"
    SENT_DATA=$(echo -ne "\x08\xfc\x71\xc7\x68\x94\xbf\x4a\x47\x2f\xb9\xd4\xd6\x46\x72\xef\x74\xf0\x6f\x30\xbc\x25\x09\x45\x4a\x05\x93\x2b\xc4\xe7\x57\x9c\x26\x54\xc4\xed\xb7\xcd\x40\x29\x2b\x2e\x28\xf3\x91\xd7\xa2\x24\x1f\x77\x93\xae\x4b\x2a\x23\xd2\x62\x17\xe6\x11\xe5\x85\x86\x53\xe6\xf7\xf5\xd3\xa1\x96\x02\x37\x2e\x06" | base64);
fi

JSON_DATA='{
        "method":"POST",
        "host":"hidden-lake-messenger",
        "path":"/push",
        "head":{
            "Hl-Messenger-Pseudonym": "Bob",
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
