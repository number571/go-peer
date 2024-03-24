N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hle.yml .
build: default-config
	$(GC) -o $(BINPATH)/hle ./cmd/hle

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_linux ./cmd/hle; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_windows.exe ./cmd/hle; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_darwin ./cmd/hle; \
	done;
run:
	./$(BINPATH)/hle
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hle $(BINPATH)/hle_*
	rm -rf hle.yml priv.key
