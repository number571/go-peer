NICKNAME=init_node
HOST=127.0.0.1
PORT=8080

.PHONY: default build run clean 
default: build run 
build: 
	cd _init && go run ../keygen
run:
	go run . $(NICKNAME) $(HOST):$(PORT)
clean:
	rm _init/priv.key _init/pub.key 
