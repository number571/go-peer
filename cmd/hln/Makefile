# Hidden Lake Network
.PHONY: default build run clean 
default: build
build:
	cp -r ../hls . 
	cp -r ../hms . 
	make build -C hls/ 
	make build -C hms/ 
	./hls/hls --init
	./hms/hms --init
run:
	(./hms/hms) & (sleep 1 && ./hls/hls)
clean:
	make clean -C ../hls
	make clean -C ../hms
	rm -rf hls hls.cfg hls.db hls.stg 
	rm -rf hms hms.cfg hms.db 
