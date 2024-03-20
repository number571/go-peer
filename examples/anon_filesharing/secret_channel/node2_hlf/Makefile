GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hlc2 &
clean:
	pkill -15 prog_hlc2 || true
	rm -rf prog_hlc2 hls.db hls.yml hlt.db hlf.stg hla_common_consumer.db 
