N=1
GC=go build
.PHONY: default build test clean
default: build
build:
	make -C ./recv build
	make -C ./send build
	make -C ./service build
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	make -C ./recv clean 
	make -C ./send clean 
	make -C ./service clean 
