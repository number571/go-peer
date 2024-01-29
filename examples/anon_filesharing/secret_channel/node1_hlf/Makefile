GC=go build
.PHONY: default run clean 
default: clean run 
run:
	cp hls_copy.yml hls.yml 
	./prog_hls1 &
	./prog_hlt1 &
	./prog_hlf1 &
	./prog_send1 8082 6060 &
	./prog_recv1 6060 8582 0 &
clean:
	pkill -15 prog_hlf1 || true
	pkill -15 prog_hls1 || true
	pkill -15 prog_hlt1 || true
	pkill -15 prog_send1 || true 
	pkill -15 prog_recv1 || true 
	rm -rf prog_hls1 hls.db hls.yml 
	rm -rf prog_hlt1 hlt.db
	rm -rf prog_hlf1 hlf.stg/loaded/*
	rm -rf prog_send1 
	rm -rf prog_recv1 common_recv.db 
