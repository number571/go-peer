N=1

PPROF_PATH=
PPROF_PORT=

CHECK_RETURN_CODE=if [ $$? != 0 ]; then exit; fi

.PHONY: default clean \
	test-prerun test-run test-coverage test-benchmark \
	pprof-run \
	git-status git-push 

default: clean test-prerun test-run

clean:
	make -C ./cmd/hidden_lake clean 
	make -C ./examples/_cmd clean


### TEST
# example run: make test-run N=10

test-prerun:
	go vet ./...;
	$(CHECK_RETURN_CODE);

	go test -race ./...;
	$(CHECK_RETURN_CODE);

test-run:
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		go test -count=1 `go list ./...`; \
		$(CHECK_RETURN_CODE); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

test-coverage:
	go test -coverprofile=test/result/coverage.out `go list ./...`
	go tool cover -html=test/result/coverage.out

test-benchmark:
	# TODO 


### PPROF
# example run: make pprof-run PPROF_PATH=hls PPROF_PORT=62109

pprof-run:
	go tool pprof -png -output test/pprof/$(PPROF_PATH)/profile.png http://localhost:$(PPROF_PORT)/debug/pprof/profile
	go tool pprof -png -output test/pprof/$(PPROF_PATH)/heap.png http://localhost:$(PPROF_PORT)/debug/pprof/heap
	go tool pprof -png -output test/pprof/$(PPROF_PATH)/goroutine.png http://localhost:$(PPROF_PORT)/debug/pprof/goroutine


### GIT

git-status:
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 
