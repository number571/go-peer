N=1hlst
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp ../../_default_cfg/hls.cfg .
	cp ../../_default_cfg/hlt.cfg .
build: default-configs
	$(GC) -o $(BINPATH)/hlst ./cmd/hlst

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlst_$${arch}_linux ./cmd/hlst; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlst_$${arch}_windows ./cmd/hlst; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlst_$${arch}_darwin ./cmd/hlst; \
	done;
run:
	./$(BINPATH)/hlst
clean:
	rm -f $(BINPATH)/hlst $(BINPATH)/hlst_*
	rm -rf hls.cfg hlt.cfg hls.db hlt.db
