GC=go build
ROOPATH=../..
BINPATH=$(ROOPATH)/bin

NICKNAME=init_node
HOST=127.0.0.1
PORT=8080

.PHONY: default genkey build run clean 
default: build run 
build:
	$(GC) -o $(BINPATH)/ma . # main 

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/ma_$${arch}_linux .; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/ma_$${arch}_windows .; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/ma_$${arch}_darwin .; \
	done;
run:
	./$(BINPATH)/ma $(NICKNAME) $(HOST):$(PORT)
genkey: 
	cd _init && go run ../keygen
clean:
	rm -f $(BINPATH)/ma $(BINPATH)/ma_*
	rm -f _init/priv.key _init/pub.key 
