GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	cp hls_copy.yml hls.yml 
run:
	./prog_hlc1 &
	./prog_producer1 :8082 localhost:6060 false &
	./prog_consumer1 localhost:6060 localhost:8582 false &
clean:
	pkill -15 prog_hlc1 || true
	pkill -15 prog_producer1 || true 
	pkill -15 prog_consumer1 || true 
	rm -rf prog_hlc1 hls.db hls.yml hlt.db hlm.db
	rm -rf prog_producer1 
	rm -rf prog_consumer1 common_consumer.db 
