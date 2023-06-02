N=1
GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin
MOBILEPATH=./cmd/m_hls

.PHONY: default all-build build run test all-clean clean \
	mobile-default mobile-build mobile-install mobile-clean \
	docker-default docker-build docker-run docker-clean 

default: build run

all-build: build mobile-build
all-clean: clean mobile-clean 

build:
	$(GC) -o $(BINPATH)/hls ./cmd/hls # main 
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hls_linux ./cmd/hls
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hls_windows.exe ./cmd/hls
run:
	./$(BINPATH)/hls
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hls $(BINPATH)/hls_linux $(BINPATH)/hls_windows.exe 
	rm -f hls.cfg hls.db

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hls"
mobile-build:
	cd ./cmd/m_hls && fyne package -os android -appID hidden_lake.service -icon icon.png
	cp ./cmd/m_hls/m_hls.apk $(BINPATH)/hls_android.apk
mobile-install:
	cd ./cmd/m_hls && fyne install -os android -appID hidden_lake.service -icon icon.png
mobile-clean:
	rm -f $(MOBILEPATH)/m_hls.apk $(BINPATH)/hls_android.apk
	rm -rf $(MOBILEPATH)/hidden_lake

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_service
docker-run:
	docker run -it --rm -v $$(pwd)/mounted:/mounted -p 9571:9571 -p 9572:9572 hl_service
docker-clean:
	docker rmi hl_service
