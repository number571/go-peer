GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hlc1 &
	./prog_send1 8082 6060 &
	./prog_recv1 6060 8582 0 &
clean:
	pkill -15 prog_hlc1 || true
	pkill -15 prog_send1 || true 
	pkill -15 prog_recv1 || true 
	rm -rf prog_hlc1 hls.db hls.yml hlt.db hlm.db
	rm -rf prog_send1 
	rm -rf prog_recv1 common_recv.db 
