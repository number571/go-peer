.PHONY: default clean request
default: clean 
clean:
	make -C ./routing clean
	make -C ./secret_channel/common clean
request:
	go run ./_request/main.go
