GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls0 &
	./prog_hlt0 &
clean:
	pkill -15 prog_hls0 || true
	pkill -15 prog_hlt0 || true
	rm -rf prog_hls0 prog_hlt0 hls.db hlt.db 
