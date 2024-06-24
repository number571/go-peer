# updates in 'test-coverage-badge' block
_COVERAGE_FLOOR=_ 

_CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi
_GO_TEST_LIST=\
	go list ./... | \
	grep -vsE '/cmd/hidden_lake/.*/cmd/' 

_TEST_UTILS_PATH=./../../test/utils
_TEST_RESULT_PATH=./_test/result
_TEST_PPROF_PATH=./_test/pprof

.PHONY: default build clean \
	test-coverage test-coverage-view test-coverage-treemap test-coverage-badge \
	pprof-run

default: build 
build:
	make -C ./service build
	make -C ./composite build
	make -C ./applications build
	make -C ./helpers build
	make -C ./adapters build
clean:
	make -C ./service clean
	make -C ./composite clean
	make -C ./applications clean
	make -C ./helpers clean
	make -C ./adapters clean

### TEST COVERAGE

test-coverage: clean
	go test -coverpkg=./... -coverprofile=$(_TEST_RESULT_PATH)/coverage.out -count=1 `$(_GO_TEST_LIST)`
	$(_CHECK_ERROR)

test-coverage-view:
	go tool cover -html=$(_TEST_RESULT_PATH)/coverage.out

test-coverage-treemap:
	go-cover-treemap -coverprofile=$(_TEST_RESULT_PATH)/coverage.out > $(_TEST_RESULT_PATH)/coverage.svg

test-coverage-badge: 
	$(eval _COVERAGE_FLOOR=go tool cover -func=$(_TEST_RESULT_PATH)/coverage.out | grep total: | grep -oP '([0-9])+(?=\.[0-9]+)')
	if [ `${_COVERAGE_FLOOR}` -lt 60 ]; then \
		cat $(_TEST_UTILS_PATH)/badge_coverage_template.svg | sed -e "s/{{.color}}/dc143c/g;s/{{.percent}}/`${_COVERAGE_FLOOR}`/g" > $(_TEST_RESULT_PATH)/badge_coverage.svg; \
	elif [ `${_COVERAGE_FLOOR}` -gt 80 ]; then \
		cat $(_TEST_UTILS_PATH)/badge_coverage_template.svg | sed -e "s/{{.color}}/97ca00/g;s/{{.percent}}/`${_COVERAGE_FLOOR}`/g" > $(_TEST_RESULT_PATH)/badge_coverage.svg; \
	else \
		cat $(_TEST_UTILS_PATH)/badge_coverage_template.svg | sed -e "s/{{.color}}/ff8c00/g;s/{{.percent}}/`${_COVERAGE_FLOOR}`/g" > $(_TEST_RESULT_PATH)/badge_coverage.svg; \
	fi

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
