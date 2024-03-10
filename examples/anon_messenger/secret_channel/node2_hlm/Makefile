GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hlc2 &
	./prog_producer2 :7072 localhost:6060 false &
	./prog_consumer2 localhost:6060 localhost:7582 false &
clean:
	pkill -15 prog_hlc2 || true
	pkill -15 prog_producer2 || true 
	pkill -15 prog_consumer2 || true 
	rm -rf prog_hlc2 hls.db hls.yml hlt.db hlm.db
	rm -rf prog_producer2 
	rm -rf prog_consumer2 common_consumer.db 
