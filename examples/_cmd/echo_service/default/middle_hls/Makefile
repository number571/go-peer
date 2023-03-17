GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls --key priv.key &
clean:
	pkill -15 prog_hls || true
	rm -rf prog_hls hls.db
