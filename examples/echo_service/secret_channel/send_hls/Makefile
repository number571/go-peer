.PHONY: default build run clean 
default: clean build run 
build:
	# nothing
run:
	./prog_hlc2 --key priv.key &
clean:
	pkill -15 prog_hlc2 || true
	rm -rf prog_hlc2 hls.db hlt.db common_consumer.db
