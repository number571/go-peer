N=1hlstm
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp ../../_default_cfg/hls.cfg .
	cp ../../_default_cfg/hlt.cfg .
	cp ../../_default_cfg/hlm.cfg .
build: default-configs
	$(GC) -o $(BINPATH)/hlstm ./cmd/hlstm

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlstm_$${arch}_linux ./cmd/hlstm; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlstm_$${arch}_windows ./cmd/hlstm; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlstm_$${arch}_darwin ./cmd/hlstm; \
	done;
run:
	./$(BINPATH)/hlstm
clean:
	rm -f $(BINPATH)/hlstm $(BINPATH)/hlstm_*
	rm -rf hls.cfg hlt.cfg hlm.cfg hlm.stg hls.db hlm.db hlt.db
