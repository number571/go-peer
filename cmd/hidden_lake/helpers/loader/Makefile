N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hll.yml .
build: default-config
	$(GC) -o $(BINPATH)/hll ./cmd/hll

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hll_$${arch}_linux ./cmd/hll; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hll_$${arch}_windows.exe ./cmd/hll; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hll_$${arch}_darwin ./cmd/hll; \
	done;
run:
	./$(BINPATH)/hll
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hll $(BINPATH)/hll_*
	rm -rf hll.yml
