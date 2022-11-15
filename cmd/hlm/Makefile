N=1
GC=go build
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o hlm *.go
run:
	./hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f hlm hlm.cfg
	rm -rf hlm.db
