N=1
TEST_PATH=./test/result
PPROF_PATH=./test/pprof
CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi

PPROF_NAME=_
PPROF_PORT=_

.PHONY: default clean \
	test-run test-race test-coverage test-benchmark \
	pprof-run \
	git-status git-push 

default: clean test-run

clean:
	make -C ./cmd clean 
	make -C ./examples clean

### TEST
# example run: make test-run N=10
# for i in {1..100}; do echo $i; go test -count=1 ./...; done;

test-run:
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		go test -cover -count=1 `go list ./...`; \
		$(CHECK_ERROR); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

test-coverage:
	go vet ./...;
	$(CHECK_ERROR);

	go test -coverprofile=$(TEST_PATH)/coverage.out `go list ./...`
	$(CHECK_ERROR);

	go test -race -cover -count=1 `go list ./...` | tee $(TEST_PATH)/result.out;

test-coverage-view:
	go tool cover -html=$(TEST_PATH)/coverage.out

### GIT

git-status:
	git add .
	git status 

git-push: test-run 
	git commit -m "update"
	git push 

### PPROF
# make pprof-run PPROF_NAME=hls PPROF_PORT=9573
# make pprof-run PPROF_NAME=hlt PPROF_PORT=9583
# make pprof-run PPROF_NAME=hlm PPROF_PORT=9593

pprof-run:
	curl 'http://localhost:$(PPROF_PORT)/debug/pprof/trace?seconds=5' > $(PPROF_PATH)/$(PPROF_NAME)/trace.out
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/threadcreate.png http://localhost:$(PPROF_PORT)/debug/pprof/threadcreate
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/profile.png http://localhost:$(PPROF_PORT)/debug/pprof/profile?seconds=5
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/heap.png http://localhost:$(PPROF_PORT)/debug/pprof/heap
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/goroutine.png http://localhost:$(PPROF_PORT)/debug/pprof/goroutine
	go tool pprof -png -output $(PPROF_PATH)/$(PPROF_NAME)/allocs.png http://localhost:$(PPROF_PORT)/debug/pprof/allocs
