GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service1 ../../_service
run:
	./prog_hlc1 --key priv.key &
	./prog_service1 &
	./prog_producer1 :8082 localhost:6060 true &
	./prog_consumer1 localhost:6060 localhost:8582 true &
clean:
	pkill -15 prog_hlc1 || true
	pkill -15 prog_service1 || true
	pkill -15 prog_consumer1 || true
	pkill -15 prog_producer1 || true
	rm -rf prog_hlc1 hls.db hlt.db
	rm -rf prog_consumer1 common_consumer.db
	rm -rf prog_producer1
	rm -rf prog_service1
