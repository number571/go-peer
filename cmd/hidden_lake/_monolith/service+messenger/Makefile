N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp ../../_default_cfg/hls.cfg .
	cp ../../_default_cfg/hlm.cfg .
build: default-configs
	$(GC) -o $(BINPATH)/hlsm ./cmd/hlsm

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlsm_$${arch}_linux ./cmd/hlsm; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlsm_$${arch}_windows ./cmd/hlsm; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlsm_$${arch}_darwin ./cmd/hlsm; \
	done;
run:
	./$(BINPATH)/hlsm
clean:
	rm -f $(BINPATH)/hlsm $(BINPATH)/hlsm_*
	rm -rf hls.cfg hlm.cfg hlm.stg hls.db hlm.db
