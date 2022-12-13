N=1
GC=go build
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o hlm ./cmd/hlm
run:
	./hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f hlm hlm.cfg hlm.stg
	rm -rf hlm.db
