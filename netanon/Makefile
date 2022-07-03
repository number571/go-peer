test:
	# NEED FIX HLS
	for i in {1..$(N)}; do \
		go clean -testcache; \
		echo $$i; \
		go test ./...; \
		if [ $$? != 0 ]; then \
			exit; \
		fi; \
	done
