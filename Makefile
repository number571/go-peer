T=update
N=1
.PHONY: default build clean test bench
default: test
push: test 
	if [ $$? != 0 ]; then \
		exit; \
	fi; \
	git add .
	git commit -m "$(T)"
	git push 
build:
	make build -C cmd/hln
clean:
	make clean -C cmd/hln
test:
	for i in {1..$(N)}; do \
		go clean -testcache; \
		echo $$i; \
		go test ./...; \
		if [ $$? != 0 ]; then \
			exit; \
		fi; \
	done
bench:
	go clean -testcache
	go test -cover -bench=. -benchmem -benchtime=$(N)x ./...
