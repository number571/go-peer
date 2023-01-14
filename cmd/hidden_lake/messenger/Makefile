N=1
GC=go build
.PHONY: default build run test clean
default: build run
build:
	make build -C ../service
	cp ../service/hls prog_hls
	$(GC) -o hlm ./cmd/hlm
run:
	./prog_hls &
	sleep 8
	./hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	pkill prog_hls || true
	rm -f prog_hls hls.cfg 
	rm -f hlm hlm.cfg hlm.stg
	rm -rf hls.db hlm.db
