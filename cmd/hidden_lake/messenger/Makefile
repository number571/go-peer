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
	$(GC) -o $(BINPATH)/hlm ./cmd/hlm
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hlm_linux ./cmd/hlm
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hlm_windows.exe ./cmd/hlm
run:
	./$(BINPATH)/hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	pkill -15 hls || true
	make clean -C ../service
	rm -f $(BINPATH)/hlm $(BINPATH)/hlm_linux $(BINPATH)/hlm_windows.exe 
	rm -rf hlm.cfg hlm.db hlm.stg

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_messenger
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9591:9591 -p 9592:9592 hl_messenger
docker-clean:
	rm -f ./_mounted/hlm.*
	docker rmi hl_messenger
