.PHONY: default build run clean 
default: clean build run 
build:
	# nothing
run:
	./prog_hlc2 --key priv.key &
	./prog_producer2 :7072 localhost:6060 true &
	./prog_consumer2 localhost:6060 localhost:7582 true &
clean:
	pkill -15 prog_hlc2 || true
	pkill -15 prog_consumer2 || true
	pkill -15 prog_producer2 || true
	rm -rf prog_hlc2 hls.db hlt.db
	rm -rf prog_consumer2 common_consumer.db
	rm -rf prog_producer2
