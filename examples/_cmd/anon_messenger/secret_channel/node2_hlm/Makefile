GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 &
	./prog_hlt2 &
	./prog_hlm2 &
	./prog_send2 7072 6060 &
	./prog_recv2 6060 7573 &
clean:
	pkill -15 prog_hlm2 || true
	sleep 3
	pkill -15 prog_hls2 || true
	pkill -15 prog_hlt2 || true
	pkill -15 prog_send2 || true 
	pkill -15 prog_recv2 || true 
	rm -rf prog_hls2 hls.db 
	rm -rf prog_hlt2 hlt.db
	rm -rf prog_hlm2 hlm.db
	rm -rf prog_send2 
	rm -rf prog_recv2 common_recv.db 
