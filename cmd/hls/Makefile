N=1
GC=go build
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o hls *.go
run:
	./hls	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f hls hls.cfg hls.stg
	rm -rf hls.db
