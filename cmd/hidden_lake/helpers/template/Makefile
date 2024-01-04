N=1
GC=go build
ROOPATH=../../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean \
	docker-default docker-build docker-run docker-clean 

default: build run
default-config: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hl_t.yml .
build: default-config
	$(GC) -o $(BINPATH)/hl_t ./cmd/hl_t

	# for arch in amd64 arm64; \
	# do \
	# 	CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hl_t_$${arch}_linux ./cmd/hl_t; \
	# 	CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hl_t_$${arch}_windows.exe ./cmd/hl_t; \
	# 	CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hl_t_$${arch}_darwin ./cmd/hl_t; \
	# done;
run:
	./$(BINPATH)/hl_t
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hl_t $(BINPATH)/hl_t_*
	rm -rf hl_t.yml

docker-default: docker-build docker-run
docker-build: docker-configs
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl__template
docker-configs:
	mkdir ./_mounted | true 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/docker/* ./_mounted
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9501:9501 hl__template
docker-clean:
	rm -rf ./_mounted/hl_t.*
	docker rmi hl__template
