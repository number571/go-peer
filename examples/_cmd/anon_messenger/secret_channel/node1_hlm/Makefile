GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls1 &
	sleep 5
	./prog_hlt1 &
	sleep 5
	./prog_hlm1 &
	sleep 5
	./prog_send1 8082 6060 &
	sleep 5
	./prog_recv1 6060 8573 &
clean:
	pkill prog_hls1 || true
	pkill prog_hlt1 || true
	pkill prog_hlm1 || true
	pkill prog_send1 || true 
	pkill prog_recv1 || true 
	rm -rf prog_hls1 hls.db
	rm -rf prog_hlt1 hlt.db
	rm -rf prog_hlm1 hlm.db
	rm -rf prog_send1 
	rm -rf prog_recv1 recv.db 
