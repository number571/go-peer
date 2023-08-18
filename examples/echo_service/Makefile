.PHONY: default clean request
default: clean 
clean:
	make -C ./default clean
	make -C ./secret_channel clean
request:
	go run ./_request/main.go
