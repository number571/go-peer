N=1
GC=go build
ROOPATH=../../..
BINPATH=$(ROOPATH)/bin
MOBILEPATH=./cmd/m_hlm

.PHONY: default all-build build run test all-clean clean \
	mobile-default mobile-build mobile-install mobile-clean \
	docker-default docker-build docker-run docker-clean 

default: build run

all-build: build mobile-build
all-clean: clean mobile-clean 

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
	rm -f hlm.cfg hlm.db hlm.stg

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hlm"
mobile-build:
	make mobile-build -C ../service
	cd ./cmd/m_hlm && fyne package -os android -appID hidden_lake.messenger -icon icon.png
	cp ./cmd/m_hlm/m_hlm.apk $(BINPATH)/hlm_android.apk
mobile-install:
	make mobile-install -C ../service
	cd ./cmd/m_hlm && fyne install -os android -appID hidden_lake.messenger -icon icon.png
mobile-clean:
	make mobile-clean -C ../service 
	rm -f $(MOBILEPATH)/m_hlm.apk $(BINPATH)/hlm_android.apk
	rm -rf $(MOBILEPATH)/hidden_lake

docker-default: docker-build docker-run
docker-build:
	docker build -f ./Dockerfile $(ROOPATH)/ --tag hl_messenger
docker-run:
	docker run -it --rm -v $$(pwd)/_mounted:/mounted -p 9591:9591 -p 9592:9592 hl_messenger
docker-clean:
	rm -f ./_mounted/hlm.*
	docker rmi hl_messenger
