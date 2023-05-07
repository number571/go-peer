N=1
GC=go build
BINPATH=../../../bin
MOBILEPATH=./cmd/m_hlt

.PHONY: default all-build build run test all-clean clean \
	mobile-default mobile-build mobile-install mobile-clean

default: build run

all-build: build mobile-build 
all-clean: clean mobile-clean 

build:
	$(GC) -o $(BINPATH)/hlt ./cmd/hlt
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hlt_linux ./cmd/hlt
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hlt_windows.exe ./cmd/hlt
run:
	./$(BINPATH)/hlt
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlt $(BINPATH)/hlt_linux $(BINPATH)/hlt_windows.exe 
	rm -f hlt.cfg hlt.db

mobile-default: mobile-build mobile-install
# fyne not supported "--sourceDir cmd/m_hlt"
mobile-build:
	cd ./cmd/m_hlt && fyne package -os android -appID hidden_lake.traffic -icon icon.png
	cp ./cmd/m_hlt/m_hlt.apk $(BINPATH)/hlt_android.apk
mobile-install:
	cd ./cmd/m_hlt && fyne install -os android -appID hidden_lake.traffic -icon icon.png
mobile-clean:
	rm -f $(MOBILEPATH)/m_hlt.apk $(BINPATH)/hlt_android.apk
	rm -rf $(MOBILEPATH)/hidden_lake
