GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hlc.yml .
build: default-configs
	$(GC) -o $(BINPATH)/hlc ./cmd/hlc

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_linux ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_windows.exe ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_darwin ./cmd/hlc; \
	done;
run:
	./$(BINPATH)/hlc
clean:
	rm -f $(BINPATH)/hlc $(BINPATH)/hlc_*
	rm -rf hl?.* priv.key
