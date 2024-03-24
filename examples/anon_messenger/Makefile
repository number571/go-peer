.PHONY: default clean request
default: clean 
clean:
	make -C ./default clean
	make -C ./group clean
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./secret_channel/chatingar clean
	make -C ./secret_channel/common clean
request:
	go run ./_request/main.go 
