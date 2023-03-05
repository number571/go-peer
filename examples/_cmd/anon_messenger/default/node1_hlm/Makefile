GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls1 &
	sleep 5
	./prog_hlm1 &
clean:
	pkill prog_hls1 || true
	pkill prog_hlm1 || true
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlm1 hlm.db
