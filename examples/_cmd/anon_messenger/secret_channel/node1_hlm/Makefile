GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls1 &
	./prog_hlt1 &
	./prog_hlm1 &
	./prog_send1 8082 6060 &
	./prog_recv1 6060 8573 &
clean:
	pkill -15 prog_hlm1 || true
	sleep 3
	pkill -15 prog_hls1 || true
	pkill -15 prog_hlt1 || true
	pkill -15 prog_send1 || true 
	pkill -15 prog_recv1 || true 
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlt1 hlt.db
	rm -rf prog_hlm1 hlm.db
	rm -rf prog_send1 
	rm -rf prog_recv1 common_recv.db 
