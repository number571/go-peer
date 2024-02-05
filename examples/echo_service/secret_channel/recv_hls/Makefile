GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service1 service/*.go
run:
	./prog_hlc1 --key priv.key &
	./prog_service1 &
	./prog_send1 8082 6060 &
	./prog_recv1 6060 8582 1 &
clean:
	pkill -15 prog_hlc1 || true
	pkill -15 prog_service1 || true
	pkill -15 prog_recv1 || true
	pkill -15 prog_send1 || true
	rm -rf prog_hlc1 hls.db hlt.db
	rm -rf prog_recv1 common_recv.db
	rm -rf prog_send1
	rm -rf prog_service1
