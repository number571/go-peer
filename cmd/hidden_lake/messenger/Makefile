N=1
GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean \
	docker-default docker-build docker-run docker-clean 

default: build run
default-config: 
	cp ../_default_cfg/hlm.cfg .
build: default-config
	$(GC) -o $(BINPATH)/hlm ./cmd/hlm
	
	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlm_$${arch}_linux ./cmd/hlm; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlm_$${arch}_windows.exe ./cmd/hlm; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlm_$${arch}_darwin ./cmd/hlm; \
	done;
run:
	./$(BINPATH)/hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlm $(BINPATH)/hlm_*
	rm -rf hlm.cfg hlm.db hlm.stg

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_messenger
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9591:9591 -p 9592:9592 hl_messenger
docker-clean:
	rm -rf ./_mounted/hlm.*
	docker rmi hl_messenger
