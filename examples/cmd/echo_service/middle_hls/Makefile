GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_hls hls/*.go
run:
	./prog_hls --key priv.key &
	sleep 5
clean:
	pkill prog_hls || true
	rm -rf prog_hls hls.db
