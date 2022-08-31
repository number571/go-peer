GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_hls hls/*.go
run:
	./prog_hls --input-key priv.key &
	sleep .5
clean:
	rm -f prog_hls 
	pkill prog_hls || true
