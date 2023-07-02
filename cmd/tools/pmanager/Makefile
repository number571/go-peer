GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build 
default: build 
build: 
	$(GC) -o $(BINPATH)/tpmanager . # main 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/tpmanager_linux .
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/tpmanager_windows.exe .
clean:
	rm -f $(BINPATH)/tpmanager $(BINPATH)/tpmanager_linux $(BINPATH)/tpmanager_windows.exe 
