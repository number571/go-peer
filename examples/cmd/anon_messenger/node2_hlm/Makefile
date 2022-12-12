GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hlm.cfg hlm/ 
	$(GC) -o hlm/prog_hlm2 hlm/*.go
run:
	./prog_hls2 --key priv.key &
	sleep 5
	(cd hlm/; ./prog_hlm2 &)
clean:
	pkill prog_hls2 || true
	pkill prog_hlm2 || true
	rm -rf prog_hls2 hls.db 
	rm -rf hlm/prog_hlm2 hlm/hlm.cfg hlm/hlm.db
