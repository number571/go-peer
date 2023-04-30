N=1
GC=go build
BINPATH=../../../bin
MOBILEPATH=./cmd/m_hls

.PHONY: default build run test clean \
	mobile-default mobile-build mobile-install mobile-clean

default: build run
build:
	$(GC) -o $(BINPATH)/hls ./cmd/hls
run:
	./$(BINPATH)/hls
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hls hls.cfg
	rm -rf hls.db

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hls"
mobile-build:
	cd ./cmd/m_hls && fyne package -os android -appID hidden_lake.service -icon icon.png
	cp ./cmd/m_hls/m_hls.apk $(BINPATH)/m_hls.apk
mobile-install:
	cd ./cmd/m_hls && fyne install -os android -appID hidden_lake.service -icon icon.png
mobile-clean:
	rm -f ./cmd/m_hls/m_hls.apk $(BINPATH)/m_hls.apk
	rm -f $(MOBILEPATH)/hls.db $(MOBILEPATH)/hls.cfg
