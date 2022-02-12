N=1
.PHONY: default clean test bench 
default: test
clean:
	go clean -testcache
test: clean
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
bench: clean
	go test -bench=. -benchmem -benchtime=$(N)x ./...
