GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build clean 
default: build 
build: 
	$(GC) -o $(BINPATH)/tkeygen . # main 
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/tkeygen_linux .
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/tkeygen_windows.exe .
clean:
	rm -f *.key 
	rm -f $(BINPATH)/tkeygen $(BINPATH)/tkeygen_linux $(BINPATH)/tkeygen_windows.exe 
