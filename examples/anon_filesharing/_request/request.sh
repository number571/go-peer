#!/bin/bash

## node1[localhost:8080](Alice) -> node2[localhost:7070](Bob)

JSON_DATA='{
        "method":"GET",
        "host":"hidden-lake-filesharer",
        "path":"/list?page=0",
        "head":{
            "Accept": "application/json"
        }
}';

JSON_DATA=${JSON_DATA//\"/\\\"} # "method" -> \"method\", ...
JSON_DATA=${JSON_DATA//[$'\t\r\n ']} # delete \t \r \n ' ' from string

PUSH_FORMAT='{
        "receiver":"Bob",
        "req_data":"'$JSON_DATA'"
}';

curl -i -X POST -H 'Accept: application/json' http://localhost:8572/api/network/request --data "${PUSH_FORMAT}";
echo 
