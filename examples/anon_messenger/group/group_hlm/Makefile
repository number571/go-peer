.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hls0 --key priv.key &
	./prog_hlm0 &
clean:
	pkill -15 prog_hlm0 || true
	pkill -15 prog_hls0 || true
	rm -rf prog_hlm0 hlm.db
	rm -rf prog_hls0 hls.db hls.yml
