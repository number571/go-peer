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

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlt_$${arch}_linux ./cmd/hlt; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlt_$${arch}_windows ./cmd/hlt; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlt_$${arch}_darwin ./cmd/hlt; \
	done;
run:
	./$(BINPATH)/hlt
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlt $(BINPATH)/hlt_*
	rm -rf hlt.cfg hlt.db

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_traffic
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9581:9581 hl_traffic
docker-clean:
	rm -rf ./_mounted/hlt.*
	docker rmi hl_traffic
