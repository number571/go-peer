N=1
TEST_PATH=./test/result
PPROF_PATH=./test/pprof
CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi

PPROF_NAME=_
PPROF_PORT=_

.PHONY: default clean \
	test-run test-coverage test-benchmark \
	pprof-run \
	git-status git-push 

default: clean test-run

clean:
	make -C ./cmd/hidden_lake clean 
	make -C ./examples/_cmd clean


### TEST
# example run: make test-run N=10
test-run:
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		go test -race -cover -count=1 `go list ./...`; \
		$(CHECK_ERROR); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

test-coverage:
	go vet ./...;
	$(CHECK_ERROR);

	go test -coverprofile=$(TEST_PATH)/coverage.out `go list ./...`
	$(CHECK_ERROR);

	go test -race -cover -count=1 `go list ./...` | tee $(TEST_PATH)/result.out;
	$(CHECK_ERROR);
	
	go tool cover -html=$(TEST_PATH)/coverage.out
	$(CHECK_ERROR);

test-benchmark:
	# TODO 


### PPROF
# example run: make pprof-run PPROF_NAME=hls PPROF_PORT=62109

pprof-run:
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/profile.png http://localhost:$(PPROF_PORT)/debug/pprof/profile
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/heap.png http://localhost:$(PPROF_PORT)/debug/pprof/heap
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/goroutine.png http://localhost:$(PPROF_PORT)/debug/pprof/goroutine


### GIT

git-status:
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 
