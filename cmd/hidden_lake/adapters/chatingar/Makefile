N=1
GC=go build
.PHONY: default build test clean
default: build
build:
	make -C ./consumer build
	make -C ./producer build
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	make -C ./consumer clean 
	make -C ./producer clean 
