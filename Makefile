T=update
N=1

.PHONY: default push test

default: test

push: test 
	if [ $$? != 0 ]; then \
		exit; \
	fi; \
	git add .
	git commit -m "$(T)"
	git push 

test:
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		go clean -testcache; \
		echo $$i; \
		go test `go list ./... | grep -v examples`; \
		if [ $$? != 0 ]; then \
			exit; \
		fi; \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";
