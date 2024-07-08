#!/bin/bash

## node1[localhost:8080](Alice) -> node2[localhost:7070](Bob)

PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"GET",
        "host":"hidden-lake-filesharer",
        "path":"/list?page=0"
    }
}';

curl -i -X POST -H 'Accept: application/json' http://localhost:8572/api/network/request --data "${PUSH_FORMAT}";
echo 
