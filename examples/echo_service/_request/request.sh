#!/bin/bash

JSON_DATA='{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
                "Accept": "application/json"
        },
        "body":"aGVsbG8sIHdvcmxkIQ=="
}';

JSON_DATA=${JSON_DATA//\"/\\\"} # "method" -> \"method\", ...
JSON_DATA=${JSON_DATA//[$'\t\r\n ']} # delete ['\t', '\r', '\n', ' '] from string

PUSH_FORMAT='{
        "receiver":"Bob",
        "req_data":"'$JSON_DATA'"
}';

d="$(date +%s)";
curl -i -X POST -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
