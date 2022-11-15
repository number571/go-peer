GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp priv.key pub.key hls.cfg hls/
	cp hlm.cfg hlm/ 
	$(GC) -o hls/prog_hls2 hls/*.go
	$(GC) -o hlm/prog_hlm2 hlm/*.go
run:
	(cd hls/; ./prog_hls2 --key priv.key &)
	sleep 5
	(cd hlm/; ./prog_hlm2 &)
clean:
	pkill prog_hls2 || true
	pkill prog_hlm2 || true
	rm -rf hls/prog_hls2 hls/hls.cfg hls/hls.db hls/priv.key hls/pub.key
	rm -rf hlm/prog_hlm2 hlm/hlm.cfg hlm/hlm.db
