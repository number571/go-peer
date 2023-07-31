N=1hlc_st
GC=go build
ROOPATH=../../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp $(ROOPATH)/cmd/hidden_lake/_default_cfg/hls.cfg .
	cp $(ROOPATH)/cmd/hidden_lake/_default_cfg/hlt.cfg .
build: default-configs
	$(GC) -o $(BINPATH)/hlc_st ./cmd/hlc_st

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_st_$${arch}_linux ./cmd/hlc_st; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_st_$${arch}_windows.exe ./cmd/hlc_st; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_st_$${arch}_darwin ./cmd/hlc_st; \
	done;
run:
	./$(BINPATH)/hlc_st
clean:
	rm -f $(BINPATH)/hlc_st $(BINPATH)/hlc_st_*
	rm -rf hls.cfg hlt.cfg hls.db hlt.db
