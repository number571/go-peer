.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hls0 --key priv.key &
clean:
	pkill -15 prog_hls0 || true
	rm -rf prog_hls0 prog_service hls.yml hls.db
