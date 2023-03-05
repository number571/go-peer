GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hls2 --key priv.key &
	sleep 5
	./prog_hlt2 &
	sleep 5
	./prog_send2 7072 6060 &
	sleep 5
	./prog_recv2 6060 7573 &
clean:
	pkill prog_hls2 || true
	pkill prog_hlt2 || true
	pkill prog_recv2 || true
	pkill prog_send2 || true
	rm -rf prog_hls2 hls.db
	rm -rf prog_hlt2 hlt.db
	rm -rf prog_recv2 recv.db
	rm -rf prog_send2
