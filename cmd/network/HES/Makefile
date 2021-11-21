CC=gcc
GC=go build
GFILES=client.go cdatabase.go cmodels.go csessions.go gconsts.go server.go sdatabase.go sconfig.go
.PHONY: default build build_all clean
default: build
build: $(GFILES)
	$(GC) client.go gconsts.go cdatabase.go cmodels.go csessions.go
	$(GC) server.go gconsts.go sdatabase.go sconfig.go
build_all: compile.c $(GFILES)
	$(CC) -o compile compile.c
	./compile
clean:
	rm -f client.db server.db server client compile client_*_* server_*_*
