N=1
GC=go build
BINPATH=../../../bin
MOBILEPATH=./cmd/m_hlt

.PHONY: default build run test all-clean clean \
	mobile-default mobile-build mobile-install mobile-clean
all-clean: clean mobile-clean 

default: build run
build:
	$(GC) -o $(BINPATH)/hlt ./cmd/hlt
run:
	./$(BINPATH)/hlt
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -rf $(BINPATH)/hlt $(BINPATH)/hlt.cfg $(BINPATH)/hlt.db

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hlt"
mobile-build:
	cd ./cmd/m_hlt && fyne package -os android -appID hidden_lake.traffic -icon icon.png
	cp ./cmd/m_hlt/m_hlt.apk $(BINPATH)/m_hlt.apk
mobile-install:
	cd ./cmd/m_hlt && fyne install -os android -appID hidden_lake.traffic -icon icon.png
mobile-clean:
	rm -f ./cmd/m_hlt/m_hlt.apk $(BINPATH)/m_hlt.apk
	rm -rf $(MOBILEPATH)/hidden_lake
