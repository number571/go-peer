N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hld.yml .
build: default-config
	$(GC) -o $(BINPATH)/hld ./cmd/hld

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hld_$${arch}_linux ./cmd/hld; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hld_$${arch}_windows.exe ./cmd/hld; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hld_$${arch}_darwin ./cmd/hld; \
	done;
run:
	./$(BINPATH)/hld
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hld $(BINPATH)/hld_*
	rm -rf hld.yml
