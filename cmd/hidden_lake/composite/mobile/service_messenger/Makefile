N=1
GC=go build
ROOPATH=../../../../..
BINPATH=$(ROOPATH)/bin
MOBILEPATH=./cmd/m_hlc_sm

.PHONY: default build run install clean

default: build 
# fyne not supported "--sourceDir cmd/m_hlc_sm"
run:
	go run $(MOBILEPATH)
build:
	cd $(MOBILEPATH) && fyne package -os android -appID hidden_lake.messenger -icon icon.png
	cp $(MOBILEPATH)/m_hlc_sm.apk $(BINPATH)/hlc_sm_arm_android.apk
install:
	cd $(MOBILEPATH) && fyne install -os android -appID hidden_lake.messenger -icon icon.png
clean:
	rm -f $(MOBILEPATH)/m_hlc_sm.apk $(BINPATH)/hlc_sm_arm_android.apk
	rm -rf ./hidden_lake/
