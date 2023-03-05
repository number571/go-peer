GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 &
	sleep 5
	./prog_hlm2 &
clean:
	pkill prog_hls2 || true
	pkill prog_hlm2 || true
	rm -rf prog_hls2 hls.db 
	rm -rf prog_hlm2 hlm.db
