N=1

.PHONY: default clean test  
default: clean test

clean:
	make -C ./cmd/hls clean
	make -C ./cmd/hlm clean
	make -C ./cmd/hms clean
	make -C ./cmd/hmc clean
	make -C ./cmd/ubc clean
	make -C ./examples/anon_messenger clean 
	make -C ./examples/echo_service clean 

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
