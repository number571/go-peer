GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.cfg hls.cfg 
run:
	./prog_hls2 &
	./prog_hlm2 &
clean:
	pkill -15 prog_hlm2 || true
	pkill -15 prog_hls2 || true
	rm -rf prog_hls2 hls.db hls.cfg
	rm -rf prog_hlm2 hlm.db
