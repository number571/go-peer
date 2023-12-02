N=1
GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean \
	docker-default docker-build docker-run docker-clean 

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hls.yml .
build: default-config
	$(GC) -o $(BINPATH)/hls ./cmd/hls # main 

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hls_$${arch}_linux ./cmd/hls; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hls_$${arch}_windows.exe ./cmd/hls; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hls_$${arch}_darwin ./cmd/hls; \
	done;
run:
	./$(BINPATH)/hls
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hls $(BINPATH)/hls_*
	rm -rf hls.yml hls.db

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_service
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9571:9571 -p 9572:9572 hl_service
docker-clean:
	rm -rf ./_mounted/hls.*
	docker rmi hl_service
