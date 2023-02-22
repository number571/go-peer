N=1
CHECK_RETURN_CODE=if [ $$? != 0 ]; then exit; fi

.PHONY: default clean \
	test-run test-coverage test-benchmark \
	git-status git-push 

default: clean test-run

clean:
	make -C ./cmd/hidden_lake clean 
	make -C ./examples/_cmd clean

test-run:
	go vet ./...;
	$(CHECK_RETURN_CODE);

	go test -race ./...;
	$(CHECK_RETURN_CODE);

	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		echo $$i; \
		go test -count=1 -coverprofile=test/coverage.out `go list ./...`; \
		$(CHECK_RETURN_CODE); \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

test-coverage:
	go test -coverprofile=test/coverage.out `go list ./...`
	go tool cover -html=test/coverage.out

test-benchmark:
	# TODO 

git-status:
	git add .
	git status 

git-push:
	git commit -m "update"
	git push 
