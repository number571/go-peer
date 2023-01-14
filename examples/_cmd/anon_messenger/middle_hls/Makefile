GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls &
	./prog_hlt &
	sleep 10
clean:
	pkill prog_hls || true
	pkill prog_hlt || true
	rm -rf prog_hls prog_hlt hlt.db hls.db
