GC=go build
.PHONY: default run clean 
default: clean run 
run:
	cp hls_copy.cfg hls.cfg 
	cp hlm_copy.stg hlm.stg 
	./prog_hls1 &
	./prog_hlm1 &
clean:
	pkill -15 prog_hlm1 || true
	pkill -15 prog_hls1 || true
	rm -rf prog_hls1 hls.db hls.cfg hlm.stg 
	rm -rf prog_hlm1 hlm.db
