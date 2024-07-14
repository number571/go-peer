N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean

default: build run
default-config: 
	# cp $(ROOPATH)/cmd/hidden_lake/_configs/hlr.yml .
build: default-config
	$(GC) -o $(BINPATH)/hlr ./cmd/hlr

	# for arch in amd64 arm64; \
	# do \
	# 	CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlr_$${arch}_linux ./cmd/hlr; \
	# 	CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlr_$${arch}_windows.exe ./cmd/hlr; \
	# 	CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlr_$${arch}_darwin ./cmd/hlr; \
	# done;
run:
	./$(BINPATH)/hlr
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlr $(BINPATH)/hlr_*
	rm -rf hlr.yml
