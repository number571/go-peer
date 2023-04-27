N=1
GC=go build
BINPATH=../../../bin

.PHONY: default build run test clean \
	mobile-default mobile-build mobile-install mobile-clean

default: build run
build:
	make build -C ../service
	$(GC) -o $(BINPATH)/hlm ./cmd/hlm
run:
	./$(BINPATH)/hls &
	./$(BINPATH)/hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	pkill -15 hls || true
	rm -f $(BINPATH)/hls $(BINPATH)/hls.cfg 
	rm -f $(BINPATH)/hlm $(BINPATH)/hlm.cfg $(BINPATH)/hlm.stg
	rm -rf $(BINPATH)/hls.db $(BINPATH)/hlm.db

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hls"
mobile-build:
	make mobile-build -C ../service
	cd ./cmd/m_hlm && fyne package -os android -appID hidden_lake.messenger -icon icon.png
	cp ./cmd/m_hlm/m_hlm.apk $(BINPATH)/m_hlm.apk
mobile-install:
	make mobile-install -C ../service
	cd ./cmd/m_hlm && fyne install -os android -appID hidden_lake.messenger -icon icon.png
mobile-clean:
	make mobile-clean -C ../service
	rm -f ./cmd/m_hlm/m_hlm.apk $(BINPATH)/m_hlm.apk
