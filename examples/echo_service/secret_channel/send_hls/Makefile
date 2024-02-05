.PHONY: default build run clean 
default: clean build run 
build:
	# nothing
run:
	./prog_hlc2 --key priv.key &
	./prog_send2 7072 6060 &
	./prog_recv2 6060 7582 1 &
clean:
	pkill -15 prog_hlc2 || true
	pkill -15 prog_recv2 || true
	pkill -15 prog_send2 || true
	rm -rf prog_hlc2 hls.db hlt.db
	rm -rf prog_recv2 common_recv.db
	rm -rf prog_send2
