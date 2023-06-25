GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build 
default: build 
build: 
	$(GC) -o $(BINPATH)/tpmanager . # main 
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/tpmanager_linux .
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/tpmanager_windows.exe .
clean:
	rm -f $(BINPATH)/tpmanager $(BINPATH)/tpmanager_linux $(BINPATH)/tpmanager_windows.exe 
