N=1

.PHONY: default status push clean test  
default: clean test

clean:
	make -C ./cmd/hidden_lake/service clean
	make -C ./cmd/hidden_lake/messenger clean
	make -C ./cmd/hidden_lake/traffic clean
	make -C ./examples/_cmd/anon_messenger clean 
	make -C ./examples/_cmd/echo_service clean 
	make -C ./examples/_cmd/traffic_keeper clean 

test:
	d=$$(date +%s); \
	for i in {1..$(N)}; do \
		go clean -testcache; \
		echo $$i; \
		go test `go list ./...`; \
		if [ $$? != 0 ]; then \
			exit; \
		fi; \
	done; \
	echo "Build took $$(($$(date +%s)-d)) seconds";

status:
	git add .
	git status 

push:
	git commit -m "update"
	git push 
