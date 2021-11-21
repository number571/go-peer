CC=gcc
GC=go build
GFILES=main.go client.go server.go funcs.go
.PHONY: default build run clean
default: build run
build: $(GFILES)
	$(GC) main.go funcs.go
	$(GC) client.go funcs.go
	$(GC) server.go
run: main 
	./main
clean:
	rm -f main client server config.json priv.key pub.key
