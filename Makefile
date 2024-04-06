N=1

PPROF_NAME=_
PPROF_PORT=_

# updates in 'test-coverage-badge' block
_COVERAGE_FLOOR=_ 

_TEST_RESULT_PATH=./test/result

_CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi
_GO_TEST_LIST=\
	go list ./... | \
	grep -v /examples/ | \
	grep -v /cmd/

.PHONY: default clean \
	lint-run test-run \
	test-coverage test-coverage-view \
	git-status git-push 

default: lint-run test-run

clean:
	make -C ./bin clean 
	make -C ./cmd clean 
	make -C ./examples clean

go-fmt-vet:
	go fmt ./...
	go vet ./...

### TEST
# example run: make test-run N=10
# for i in {1..100}; do echo $i; go test -race -shuffle=on -count=1 ./...; done;

test-run: clean go-fmt-vet
	$(_CHECK_ERROR);
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		# recommended to add an option -shuffle=on if [go version >= 1.17]; \
		go test -race -cover -count=1 ./...; \
		$(_CHECK_ERROR); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

### TEST COVERAGE

test-coverage: clean go-fmt-vet
	make test-coverage -C cmd/hidden_lake/
	go test -coverpkg=./... -coverprofile=$(_TEST_RESULT_PATH)/coverage.out -count=1 `$(_GO_TEST_LIST)`
	$(_CHECK_ERROR)

test-coverage-view:
	make test-coverage-view -C cmd/hidden_lake/
	go tool cover -html=$(_TEST_RESULT_PATH)/coverage.out

test-coverage-badge: 
	make test-coverage-badge -C cmd/hidden_lake/
	$(eval _COVERAGE_FLOOR=go tool cover -func=$(_TEST_RESULT_PATH)/coverage.out | grep total: | grep -oP '([0-9])+(?=\.[0-9]+)')
	if [ `${_COVERAGE_FLOOR}` -lt 60 ]; then \
		curl "https://img.shields.io/badge/coverage-`${_COVERAGE_FLOOR}`%25-crimson" > $(_TEST_RESULT_PATH)/badge.svg; \
	elif [ `${_COVERAGE_FLOOR}` -gt 80 ]; then \
		curl "https://img.shields.io/badge/coverage-`${_COVERAGE_FLOOR}`%25-green" > $(_TEST_RESULT_PATH)/badge.svg; \
	else \
		curl "https://img.shields.io/badge/coverage-`${_COVERAGE_FLOOR}`%25-darkorange" > $(_TEST_RESULT_PATH)/badge.svg; \
	fi

### LINT

lint-run:
	golangci-lint run -E "gas,unconvert,gosimple,goconst,gocyclo,goerr113,ineffassign,unparam,unused,bodyclose,noctx,perfsprint,prealloc,gocritic,govet,revive,staticcheck,errcheck,errorlint,nestif,maintidx"

### GIT

git-status: lint-run test-coverage test-coverage-badge
	go fmt ./...
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 
