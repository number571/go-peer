N=1
GC=go build
BINPATH=../../../bin
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o $(BINPATH)/hls ./cmd/hls
run:
	./$(BINPATH)/hls
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hls hls.cfg
	rm -rf hls.db
