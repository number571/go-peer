GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls --key priv.key &
	sleep 5
clean:
	pkill prog_hls || true
	rm -rf prog_hls hls.db
