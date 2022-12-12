GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hlm.cfg hlm/ 
	$(GC) -o hlm/prog_hlm1 hlm/*.go
run:
	./prog_hls1 --key priv.key &
	sleep 5
	(cd hlm/; ./prog_hlm1 &)
clean:
	pkill prog_hls1 || true
	pkill prog_hlm1 || true
	rm -rf prog_hls1 hls/hls.db
	rm -rf hlm/prog_hlm1 hlm/hlm.cfg hlm/hlm.db
