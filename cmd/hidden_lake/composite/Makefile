GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin

.PHONY: default configs build run clean

default: build run
default-configs: 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/hlc.yml .
build: default-configs
	$(GC) -o $(BINPATH)/hlc ./cmd/hlc

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_linux ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_windows.exe ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_darwin ./cmd/hlc; \
	done;
run:
	./$(BINPATH)/hlc
clean:
	rm -f $(BINPATH)/hlc $(BINPATH)/hlc_*
	rm -rf hl?.* priv.key

docker-default: docker-build docker-run
docker-build: docker-configs
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_composite
docker-configs:
	mkdir ./_mounted | true 
	cp $(ROOPATH)/cmd/hidden_lake/_configs/docker/* ./_mounted
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted \
		-p 9571:9571 \
		-p 9572:9572 \
		hl_composite
docker-clean:
	rm -rf ./_mounted/hl?.* priv.key
	docker rmi hl_composite
