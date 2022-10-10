GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_hls3 hls/*.go
	$(GC) -o prog_hlm3 hlm/*.go
run:
	./prog_hls3 --input-key priv.key &
	sleep .5
	./prog_hlm3
clean:
	pkill prog_hls3 || true
	pkill prog_hlm3 || true
	rm -rf prog_hls3 prog_hlm3 hls.db hlm.db
	