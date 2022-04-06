T=1
B=0
.PHONY: default build clean test  
default: build
build:
	make build -C cmd/hls
	make build -C cmd/hms
clean:
	make clean -C cmd/hls
	make clean -C cmd/hms
test:
	if [ $(B) == 0 ]; then \
		for i in {1..$(T)}; do go clean -testcache; echo $$i; go test ./...; done \
	else \
		go clean -testcache; go test -bench=. -benchmem -benchtime=$(B)x ./...; \
	fi
