GC=go build
N=4096

ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default genkey build clean 
default: build
build: 
	$(GC) -o $(BINPATH)/tencryptor . # main 
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/tencryptor_linux .
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/tencryptor_windows.exe .
clean:
	rm -f *.msg *.key 
	rm -f $(BINPATH)/tencryptor $(BINPATH)/tencryptor_linux $(BINPATH)/tencryptor_windows.exe 
genkey:
	go run ../keygen $(N)
