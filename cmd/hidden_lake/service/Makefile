N=1
GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default all-build build run test all-clean clean \
	docker-default docker-build docker-run docker-clean 

default: build run

all-build: build
all-clean: clean

build:
	$(GC) -o $(BINPATH)/hls ./cmd/hls # main 
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hls_linux ./cmd/hls
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hls_windows.exe ./cmd/hls
run:
	./$(BINPATH)/hls
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hls $(BINPATH)/hls_linux $(BINPATH)/hls_windows.exe 
	rm -rf hls.cfg hls.db

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_service
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9571:9571 -p 9572:9572 hl_service
docker-clean:
	rm -rf ./_mounted/hls.*
	docker rmi hl_service
