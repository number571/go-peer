GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hlc1 &
clean:
	pkill -15 prog_hlc1 || true
	rm -rf prog_hlc1 hls.db hls.yml hlt.db hlm.db common_consumer.db 
