GC=go build
.PHONY: default run clean 
default: clean run 
run:
	cp hls_copy.yml hls.yml 
	./prog_hlc2 &
	./prog_send2 7072 6060 &
	./prog_recv2 6060 7582 0 &
clean:
	pkill -15 prog_hlc2 || true
	pkill -15 prog_send2 || true 
	pkill -15 prog_recv2 || true 
	rm -rf prog_hlc2 hls.db hls.yml hlt.db hlm.db
	rm -rf prog_send2 
	rm -rf prog_recv2 common_recv.db 
