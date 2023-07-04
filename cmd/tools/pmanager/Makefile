GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build 
default: build 
build: 
	$(GC) -o $(BINPATH)/tpmanager . # main 
	
	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/tpmanager_$${arch}_linux .; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/tpmanager_$${arch}_windows .; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/tpmanager_$${arch}_darwin .; \
	done;
clean:
	rm -f $(BINPATH)/tpmanager $(BINPATH)/tpmanager_linux $(BINPATH)/tpmanager_windows.exe 
