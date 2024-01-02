N=1

PPROF_NAME=_
PPROF_PORT=_

# updates in 'test-coverage-badge' block
_COVERAGE_RAW=_ 
_COVERAGE_VAR=_

_TEST_RESULT_PATH=./test/result
_TEST_PPROF_PATH=./test/pprof

_CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi

.PHONY: default clean \
	test-run test-coverage test-coverage-view \
	pprof-run \
	git-status git-push 

default: test-run

clean:
	make -C ./cmd clean 
	make -C ./examples clean

### TEST
# example run: make test-run N=10
# for i in {1..100}; do echo $i; go test -count=1 ./...; done;

test-run: clean
	go vet ./...;
	$(_CHECK_ERROR);
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		# recommended to add an option -shuffle=on if [go version >= 1.17]; \
		go test -race -cover -count=1 ./...; \
		$(_CHECK_ERROR); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

test-coverage: clean
	go fmt ./...; \
	go test -coverpkg=./... -coverprofile=$(_TEST_RESULT_PATH)/coverage.out -count=1 ./...; \
	$(_CHECK_ERROR); \

test-coverage-view:
	go tool cover -html=$(_TEST_RESULT_PATH)/coverage.out

test-coverage-badge: 
	$(eval _COVERAGE_RAW=go tool cover -func=$(_TEST_RESULT_PATH)/coverage.out | grep total: | grep -Eo '[0-9]+\.[0-9]+')
	$(eval _COVERAGE_VAR := $(shell echo "`${_COVERAGE_RAW}`/1" | bc))
	if [ $(_COVERAGE_VAR) -lt 60 ]; then \
		curl "https://img.shields.io/badge/coverage-$(_COVERAGE_VAR)%25-red" > $(_TEST_RESULT_PATH)/badge.svg; \
	elif [ $(_COVERAGE_VAR) -gt 80 ]; then \
		curl "https://img.shields.io/badge/coverage-$(_COVERAGE_VAR)%25-green" > $(_TEST_RESULT_PATH)/badge.svg; \
	else \
		curl "https://img.shields.io/badge/coverage-$(_COVERAGE_VAR)%25-orange" > $(_TEST_RESULT_PATH)/badge.svg; \
	fi

### GIT

git-status: test-coverage test-coverage-badge
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 

### PPROF
# make pprof-run PPROF_NAME=hls PPROF_PORT=9573
# make pprof-run PPROF_NAME=hlt PPROF_PORT=9583
# make pprof-run PPROF_NAME=hlm PPROF_PORT=9593

pprof-run:
	go tool pprof -png -output $(_TEST_PPROF_PATH)/$(PPROF_NAME)/threadcreate.png http://localhost:$(PPROF_PORT)/debug/pprof/threadcreate
	go tool pprof -png -output $(_TEST_PPROF_PATH)/$(PPROF_NAME)/profile.png http://localhost:$(PPROF_PORT)/debug/pprof/profile?seconds=5
	go tool pprof -png -output $(_TEST_PPROF_PATH)/$(PPROF_NAME)/heap.png http://localhost:$(PPROF_PORT)/debug/pprof/heap
	go tool pprof -png -output $(_TEST_PPROF_PATH)/$(PPROF_NAME)/goroutine.png http://localhost:$(PPROF_PORT)/debug/pprof/goroutine
	go tool pprof -png -output $(_TEST_PPROF_PATH)/$(PPROF_NAME)/allocs.png http://localhost:$(PPROF_PORT)/debug/pprof/allocs
