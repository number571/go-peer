GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp priv.key pub.key hls.cfg hls/
	cp hlm.cfg hlm/ 
	$(GC) -o hls/prog_hls1 hls/*.go
	$(GC) -o hlm/prog_hlm1 hlm/*.go
run:
	(cd hls/; ./prog_hls1 --key priv.key &)
	sleep 5
	(cd hlm/; ./prog_hlm1 &)
clean:
	pkill prog_hls1 || true
	pkill prog_hlm1 || true
	rm -rf hls/prog_hls1 hls/hls.cfg hls/hls.db hls/priv.key hls/pub.key
	rm -rf hlm/prog_hlm1 hlm/hlm.cfg hlm/hlm.db
