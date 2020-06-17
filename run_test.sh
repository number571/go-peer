#!/bin/sh

code=0

i=1
count=5

while [ "$i" -le "$count" ]
do 
	echo "\n\n\n===== [ $i ] =====\n\n\n"
	go test -v -cover -count 10
	code="$?"
	if [ "$code" != 0 ]
	then
		break
	fi
	i=$((i+1))
done

i=$((i-1))
echo "\n= Iter: $i;\n= Code: $code;\n"
