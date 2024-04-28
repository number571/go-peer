#!/bin/bash

PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
                "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
    }
}';

d="$(date +%s)";
curl -i -X POST -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
