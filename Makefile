N=1

# updates in 'test-coverage-badge' block
_COVERAGE_FLOOR=_ 

# updates in 'git-code-lines' block
_CODE_LINES_FLOOR=_ 

_TEST_UTILS_PATH=./test/utils
_TEST_RESULT_PATH=./test/result

_CHECK_ERROR=if [ $$? != 0 ]; then exit 1; fi
_GO_TEST_LIST=\
	go list ./... | \
	grep -v /cmd/tools/ | \
	grep -v /examples/ | \
	grep -v /vendor/

.PHONY: default clean go-fmt-vet \
	lint-run test-run \
	test-coverage test-coverage-view test-coverage-treemap test-coverage-badge \
	git-status git-push \
	install-deps

default: lint-run test-run

clean:
	make -C ./bin clean 
	make -C ./cmd clean 

go-fmt-vet:
	go fmt ./...
	go vet ./...

### INSTALL

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.0
	go install github.com/nikolaydubina/go-cover-treemap@v1.4.2

### LINT

lint-run: clean go-fmt-vet
	golangci-lint run -E "gosec,unconvert,gosimple,goconst,gocyclo,err113,ineffassign,unparam,unused,bodyclose,noctx,perfsprint,prealloc,gocritic,govet,revive,staticcheck,errcheck,errorlint,nestif,maintidx"

### TEST
# example run: make test-run N=10
# for i in {1..100}; do echo $i; go test -race -shuffle=on -count=1 ./...; done;

test-run:
	$(_CHECK_ERROR);
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		go test -race -cover -count=1 -shuffle=on `$(_GO_TEST_LIST)`; \
		$(_CHECK_ERROR); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

### TEST COVERAGE

test-coverage:
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

### GIT

git-code-lines:
	$(eval _CODE_LINES_FLOOR=git ls-files | grep -v "vendor" | grep ".go" | xargs wc -l | grep total | grep -oP '([0-9])+')
	cat $(_TEST_UTILS_PATH)/badge_codelines_template.svg | sed -e "s/{{.color}}/4682b4/g;s/{{.code_lines}}/`${_CODE_LINES_FLOOR}`/g" > $(_TEST_RESULT_PATH)/badge_codelines.svg; \

git-status: lint-run test-coverage test-coverage-treemap test-coverage-badge git-code-lines
	go fmt ./...
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 
