GC=go build
.PHONY: default build run clean 
default: clean build run 
build:
	$(GC) -o prog_service1 ../../_service
run:
	./prog_hlc1 --key priv.key &
	./prog_service1 &
clean:
	pkill -15 prog_hlc1 || true
	pkill -15 prog_service1 || true
	rm -rf prog_hlc1 hls.db hlt.db common_consumer.db
	rm -rf prog_service1
