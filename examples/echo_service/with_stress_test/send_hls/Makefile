GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 --key priv.key &
clean:
	pkill -15 prog_hls2 || true
	rm -rf prog_hls2 hls.db
