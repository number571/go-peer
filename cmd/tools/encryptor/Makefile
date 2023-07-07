GC=go build
N=4096

ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default genkey build clean 
default: build
build: 
	$(GC) -o $(BINPATH)/tencryptor . # main 

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/tencryptor_$${arch}_linux .; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/tencryptor_$${arch}_windows .; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/tencryptor_$${arch}_darwin .; \
	done;
clean:
	rm -f *.msg *.key 
	rm -f $(BINPATH)/tencryptor $(BINPATH)/tencryptor_*
genkey:
	go run ../keygen $(N)
