GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service ../../_service
run:
	./prog_hls1 --key priv.key &
	./prog_service &
clean:
	pkill -15 prog_hls1 || true
	pkill -15 prog_service || true
	rm -rf prog_hls1 prog_service hls.db
