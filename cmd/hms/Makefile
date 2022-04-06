T=1
GC=go build
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o hms *.go
run:
	./hms
test:
	for i in {1..$(T)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -rf hms hms.cfg hms.db
