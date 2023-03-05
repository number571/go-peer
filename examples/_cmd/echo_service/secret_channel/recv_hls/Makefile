GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service1 service/*.go
run:
	./prog_hls1 --key priv.key &
	sleep 5
	./prog_hlt1 &
	sleep 5
	./prog_service1 &
	sleep 5
	./prog_send1 8082 6060 &
	sleep 5
	./prog_recv1 6060 8573 &
clean:
	pkill prog_hls1 || true
	pkill prog_hlt1 || true
	pkill prog_service1 || true
	pkill prog_recv1 || true
	pkill prog_send1 || true
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlt1 hlt.db
	rm -rf prog_recv1 recv.db
	rm -rf prog_send1
	rm -rf prog_service1
