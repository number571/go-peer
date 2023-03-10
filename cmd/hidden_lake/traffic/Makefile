N=1
GC=go build
BINPATH=../../../bin
.PHONY: default build run test clean
default: build run
build:
	$(GC) -o $(BINPATH)/hlt ./cmd/hlt
run:
	./$(BINPATH)/hlt
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -rf $(BINPATH)/hlt $(BINPATH)/hlt.cfg $(BINPATH)/hlt.db
