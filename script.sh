#!/bin/sh

i=1
code=0
count=150

while [ "$i" -le "$count" ]
do 
	echo "\n\n\n===== [ $i ] =====\n\n\n"
	go test -v -cover
	code="$?"
	if [ "$code" != 0 ]
	then
		break
	fi
	i=$((i+1))
done

i=$((i-1))
echo "\n= Iter: $i;\n= Code: $code;\n"
