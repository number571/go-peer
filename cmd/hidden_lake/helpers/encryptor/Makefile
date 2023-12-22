N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean \
	docker-default docker-build docker-run docker-clean 

default: build run
default-config: 
	cp ../../_configs/hle.yml .
build: default-config
	$(GC) -o $(BINPATH)/hle ./cmd/hle

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_linux ./cmd/hle; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_windows.exe ./cmd/hle; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hle_$${arch}_darwin ./cmd/hle; \
	done;
run:
	./$(BINPATH)/hle
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hle $(BINPATH)/hle_*
	rm -rf hle.yml

docker-default: docker-build docker-run
docker-build: docker-configs
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_encryptor
docker-configs:
	mkdir ./_mounted | true 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/docker/* ./_mounted
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9551:9551 hl_encryptor
docker-clean:
	rm -rf ./_mounted/hle.*
	docker rmi hl_encryptor
