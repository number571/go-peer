GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build clean
default: build 
build: 
	$(GC) -o $(BINPATH)/tstorage . # main 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/tstorage_linux .
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/tstorage_windows.exe .
clean:
	rm -f *.stg 
	rm -f $(BINPATH)/tstorage $(BINPATH)/tstorage_linux $(BINPATH)/tstorage_windows.exe 
