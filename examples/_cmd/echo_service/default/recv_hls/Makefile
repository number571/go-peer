GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service service/*.go
run:
	./prog_hls1 --key priv.key &
	sleep 5
	./prog_service &
clean:
	pkill prog_hls1 || true
	pkill prog_service || true
	rm -rf prog_hls1 prog_service hls.db
