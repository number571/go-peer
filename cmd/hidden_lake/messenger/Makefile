N=1
GC=go build
BINPATH=../../../bin
.PHONY: default build run test clean
default: build run
build:
	make build -C ../service
	$(GC) -o $(BINPATH)/hlm ./cmd/hlm
run:
	./$(BINPATH)/hls &
	sleep 8
	./$(BINPATH)/hlm	
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	pkill hls || true
	rm -f $(BINPATH)/hls $(BINPATH)/hls.cfg 
	rm -f $(BINPATH)/hlm $(BINPATH)/hlm.cfg $(BINPATH)/hlm.stg
	rm -rf $(BINPATH)/hls.db $(BINPATH)/hlm.db
