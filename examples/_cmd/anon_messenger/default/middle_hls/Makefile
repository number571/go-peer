GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls &
	sleep 5
	./prog_hlt &
clean:
	pkill prog_hls || true
	pkill prog_hlt || true
	rm -rf prog_hls prog_hlt hls.db hlt.db 
