N=1
GC=go build
BINPATH=../../../bin
MOBILEPATH=./cmd/m_hlm

.PHONY: default all-build build run test all-clean clean \
	mobile-default mobile-build mobile-install mobile-clean

default: build run

all-build: build mobile-build
all-clean: clean mobile-clean 

build:
	make build -C ../service
	$(GC) -o $(BINPATH)/hlm ./cmd/hlm
	GOARCH=amd64 GOOS=linux $(GC) -o $(BINPATH)/hlm_linux ./cmd/hlm
	GOARCH=amd64 GOOS=windows $(GC) -o $(BINPATH)/hlm_windows.exe ./cmd/hlm
run:
	./$(BINPATH)/hls &
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
	cd ./cmd/m_hlm && fyne package -os android -appID hidden_lake.messenger -icon icon.png
	cp ./cmd/m_hlm/m_hlm.apk $(BINPATH)/hlm_android.apk
mobile-install:
	cd ./cmd/m_hlm && fyne install -os android -appID hidden_lake.messenger -icon icon.png
mobile-clean:
	rm -f $(MOBILEPATH)/m_hlm.apk $(BINPATH)/hlm_android.apk
	rm -rf $(MOBILEPATH)/hidden_lake
