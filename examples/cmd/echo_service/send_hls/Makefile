GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_hls2 hls/*.go
run:
	./prog_hls2 --key priv.key &
	sleep 5
clean:
	pkill prog_hls2 || true
	rm -rf prog_hls2 hls.db
