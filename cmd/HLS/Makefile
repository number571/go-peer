CC=gcc
GC=go build
GFILES=main.go client.go server.go funcs.go
.PHONY: default build clean
default: build
build: $(GFILES)
	$(GC) main.go funcs.go
	$(GC) client.go funcs.go
	$(GC) server.go
clean:
	rm -f main client server services.json connects.json priv.key pub.key
