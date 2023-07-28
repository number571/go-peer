#!/bin/bash

# Works only if users are logged in to the account!
# node2[localhost:7070] -> node1[localhost:8080]

# body[AWhlbGxvLCB3b3JsZCE=] -> byte(0x01) || bytes("hello, world!")
# body[AmV4YW1wbGUudHh0AmhlbGxvLCB3b3JsZCE=] -> byte(0x02) || bytes("example.txt") || byte(0x02) || bytes("hello, world!")

# byte(0x01) -> text format
# byte(0x02) -> file format

str2hex() {
    local str=${1:-""}
    local fmt="%02X"
    local chr
    local -i i
    for i in `seq 0 $((${#str}-1))`; do
        chr=${str:i:1}
        printf "${fmt}" "'${chr}"
    done
}

JSON_DATA='{
        "method":"POST",
        "host":"go-peer/hidden-lake-messenger",
        "path":"/push",
        "head":{
            "Accept": "application/json"
        },
        "body":"AWhlbGxvLCB3b3JsZCE="
}';

PUSH_FORMAT='{
        "receiver":"PubKey(go-peer/rsa){3082020A0282020100C27726C08BD1A3409D99947623437AAACB3AA13929D5B65437E5B79A4E79B092A3F4C3BEABA6D7D3BA4F18357B7BE5B440F00DC66CF792D0091148A7AF8801D39D8916621A9AB52C32C4F2BAF9535FB3911392F57F455DA866760D6D1AB1F7967C5E9E007B1FF1E6854E973B724507483310981E878B6D906613A94D903F70D56BFD09596E21F40843C58C4A37142CD2320F7F466F6DC56654A0D1839681C5689D3067A10370493898E49FFBB5583E2BFB7C4712DFEFDAE6D5430E7003BFB6F4EA4CAA4C41A255D7986DFACDB4DC462912E06490B87063D4A72DF26BC705B579D50985A5E9D9E2C93CC84FF7235005DF8F689933236D9B6C54014F8ADF571559213C851CDF1BF516CB1772480289144525E944BE0397D1A8C6E58731F9E984DC5FC02D39646B6C9E9E1AA5CC76344C92F15A1E1A8632EE88DDEA5FF9E0F8DE624ACFE62C552EBD874D6379FA5CB8C349224621188FECAEF4C0FABD9513E078DC4569E6F35D64E8A0E8C8B72834F828AAB299F4F1257177B17E07C9213827088467F58101CA38EDA1DCC864EE7B3DD46AF788BEB402C001C3477FF1786C7CBF58F9FB195C7AD05D8E73081176B503DE6FE5CAC85CC39017F5E6F6EF1EB749CFAE4F277EDCA97072A1A173478984E3CF6F2A18B30368199C122BD3F926859CF3CD17ED026A59CB266EF900C107263B677C75A74A0C3716E28AFCAA51C83FE8E5FD0203010001}",
        "hex_data":"'$(str2hex "$JSON_DATA")'"
}';

while true
do
    curl -i -X PUT -H 'Accept: application/json' http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
    echo && echo && sleep 5 # seconds, queue period
done
