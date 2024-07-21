#!/bin/bash

BASE64_BODY="$(\
    echo -n 'bash[@remoter-separator]-c[@remoter-separator]echo 'hello, world' > file.txt && cat file.txt' | \
    base64 -w 0 \
)";
PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"POST",
		"host":"hidden-lake-remoter",
		"path":"/exec",
        "head":{
            "Hl-Remoter-Password": "DpxJFjAlrs4HOWga0wk14mZqQSBo9DxK"
        },
        "body":"'${BASE64_BODY}'"
    }
}';

d="$(date +%s)";
curl -i -X POST -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
