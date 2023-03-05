N=1
GC=go build
.PHONY: default build test clean
default: build
build:
	make -C recv 
	make -C send 
	make -C service 
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -rf send/send 
	rm -rf recv/recv recv/recv.db 
	rm -rf service/service service/service.db 
