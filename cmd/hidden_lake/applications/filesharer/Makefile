N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean \
	docker-default docker-build docker-run docker-clean 

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hlf.yml .
build: default-config
	$(GC) -o $(BINPATH)/hlf ./cmd/hlf
	
	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlf_$${arch}_linux ./cmd/hlf; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlf_$${arch}_windows.exe ./cmd/hlf; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlf_$${arch}_darwin ./cmd/hlf; \
	done;
run:
	./$(BINPATH)/hlf	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlm $(BINPATH)/hlf_*
	rm -rf hlf.yml hlf.stg 

docker-default: docker-build docker-run
docker-build: docker-configs
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_filesharer
docker-configs:
	mkdir ./_mounted | true 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/docker/* ./_mounted
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9541:9541 -p 9542:9542 hl_filesharer
docker-clean:
	rm -rf ./_mounted/hlf.*
	docker rmi hl_filesharer
