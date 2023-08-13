GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build clean
default: build 
build: 
	$(GC) -o $(BINPATH)/tstorage . # main 

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/tstorage_$${arch}_linux .; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/tstorage_$${arch}_windows.exe .; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/tstorage_$${arch}_darwin .; \
	done;
clean:
	rm -f *.stg 
	rm -f tstorage $(BINPATH)/tstorage $(BINPATH)/tstorage_*
