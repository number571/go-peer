GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service1 service/*.go
run:
	./prog_hls1 --key priv.key &
	./prog_hlt1 &
	./prog_service1 &
	./prog_send1 8082 6060 &
	./prog_recv1 6060 8572 &
clean:
	pkill -15 prog_hls1 || true
	pkill -15 prog_hlt1 || true
	pkill -15 prog_service1 || true
	pkill -15 prog_recv1 || true
	pkill -15 prog_send1 || true
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlt1 hlt.db
	rm -rf prog_recv1 common_recv.db
	rm -rf prog_send1
	rm -rf prog_service1
