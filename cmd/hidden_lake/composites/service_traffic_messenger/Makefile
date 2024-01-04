N=1hlc_stm
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hls.yml .
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hlt.yml .
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hlm.yml .
build: default-configs
	$(GC) -o $(BINPATH)/hlc_stm ./cmd/hlc_stm

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_stm_$${arch}_linux ./cmd/hlc_stm; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_stm_$${arch}_windows.exe ./cmd/hlc_stm; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_stm_$${arch}_darwin ./cmd/hlc_stm; \
	done;
run:
	./$(BINPATH)/hlc_stm
clean:
	rm -f $(BINPATH)/hlc_stm $(BINPATH)/hlc_stm_*
	rm -rf hls.yml hlt.yml hlm.yml hls.db hlm.db hlt.db
