GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hls3 &
	./prog_hlm3 &
clean:
	pkill -15 prog_hlm3 || true
	pkill -15 prog_hls3 || true
	rm -rf prog_hls3 hls.db hls.yml
	rm -rf prog_hlm3 hlm.db
