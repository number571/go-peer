N=1
GC=go build
ROOPATH=../../../../..
BINPATH=$(ROOPATH)/bin
MOBILEPATH=./cmd/m_hlc_stm

.PHONY: default build run install clean

default: build 
run:
	go run $(MOBILEPATH)
build:
	cd $(MOBILEPATH) && fyne package --target android -appID hidden_lake.service_traffic_messenger --icon icon.png --release --name m_hlc_stm; \

	for arch in amd64 arm64; \
	do \
		fyne package --sourceDir $(MOBILEPATH) --target android/$${arch} -appID hidden_lake.service_traffic_messenger --icon icon.png --release --name hlc_stm_$${arch}_android; \
	done;
	for arch in amd64 arm64; \
	do \
		cp $(MOBILEPATH)/hlc_stm_$${arch}_android.apk $(BINPATH)/hlc_stm_$${arch}_android.apk; \
	done;
install:
	fyne install --sourceDir $(MOBILEPATH) --target android/arm64 -appID hidden_lake.service_traffic_messenger --icon icon.png --release
clean:
	rm -f $(MOBILEPATH)/m_hlc_stm.apk $(BINPATH)/hlc_stm_*.apk
	rm -rf ./hidden_lake/
