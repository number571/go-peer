.PHONY: default run clean 
default: clean run 
run:
	./prog_hls0 --key priv.key &
clean:
	pkill -15 prog_hls0 || true
	rm -rf prog_hls0 prog_service hls.db
