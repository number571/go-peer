GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls1 &
	./prog_hlm1 &
clean:
	pkill -15 prog_hlm1 || true
	pkill -15 prog_hls1 || true
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlm1 hlm.db
