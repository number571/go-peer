GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 --key priv.key &
	sleep 5
clean:
	pkill prog_hls2 || true
	rm -rf prog_hls2 hls.db
