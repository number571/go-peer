GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp 66eecc602041c455_image.jpg.tmp hlf.stg/loaded
	cp hls_copy.yml hls.yml 
run:
	./prog_hls1 &
	./prog_hlf1 &
clean:
	pkill -15 prog_hlf1 || true
	pkill -15 prog_hls1 || true
	rm -rf prog_hls1 hls.db hls.yml
	rm -rf prog_hlf1 hlf.stg/loaded/*
