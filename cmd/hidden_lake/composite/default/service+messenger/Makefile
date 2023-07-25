N=1
GC=go build
ROOPATH=../../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp ../../_default_cfg/hls.cfg .
	cp ../../_default_cfg/hlm.cfg .
build: default-configs
	$(GC) -o $(BINPATH)/hlc_sm ./cmd/hlc_sm

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_sm_$${arch}_linux ./cmd/hlc_sm; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_sm_$${arch}_windows ./cmd/hlc_sm; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_sm_$${arch}_darwin ./cmd/hlc_sm; \
	done;
run:
	./$(BINPATH)/hlc_sm
clean:
	rm -f $(BINPATH)/hlc_sm $(BINPATH)/hlc_sm_*
	rm -rf hls.cfg hlm.cfg hlm.stg hls.db hlm.db
