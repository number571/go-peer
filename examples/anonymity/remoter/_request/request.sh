#!/bin/bash

# bash[@remoter-separator]-c[@remoter-separator]echo 'hello, world' > file.txt && cat file.txt
PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"POST",
		"host":"hidden-lake-remoter",
		"path":"/exec",
        "head":{
            "Hl-Remoter-Password": "DpxJFjAlrs4HOWga0wk14mZqQSBo9DxK"
        },
        "body":"YmFzaFtAcmVtb3Rlci1zZXBhcmF0b3JdLWNbQHJlbW90ZXItc2VwYXJhdG9yXWVjaG8gJ2hlbGxvLCB3b3JsZCcgPiBmaWxlLnR4dCAmJiBjYXQgZmlsZS50eHQ="
    }
}';

d="$(date +%s)";
curl -i -X POST -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
