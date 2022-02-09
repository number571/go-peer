N=1
.PHONY: default clean test bench 
default: clean test
clean:
	go clean -testcache
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
bench:
	go test -bench=. -benchmem -benchtime=$(N)x ./...
