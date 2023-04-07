GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 &
	./prog_hlm2 &
clean:
	pkill -15 prog_hlm2 || true
	pkill -15 prog_hls2 || true
	rm -rf prog_hls2 hls.db 
	rm -rf prog_hlm2 hlm.db
