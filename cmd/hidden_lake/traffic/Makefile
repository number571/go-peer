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
	$(GC) -o $(BINPATH)/hlt ./cmd/hlt
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hlt_linux ./cmd/hlt
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hlt_windows.exe ./cmd/hlt
run:
	./$(BINPATH)/hlt
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlt $(BINPATH)/hlt_linux $(BINPATH)/hlt_windows.exe 
	rm -rf hlt.cfg hlt.db

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_traffic
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9581:9581 hl_traffic
docker-clean:
	rm -rf ./_mounted/hlt.*
	docker rmi hl_traffic
