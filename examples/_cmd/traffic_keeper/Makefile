.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
build: 
	# MAKEFILE BUILD
	make build -C ../../../cmd/hidden_lake/traffic
	# COPY HLS
	cp -r ../../../cmd/hidden_lake/traffic/hlt ./service_hlt/prog_hlt
run: 
	# MAKEFILE RUN
	make run -C ./service_hlt
clean:
	# MAKEFILE CLEAN 
	make clean -C ./service_hlt
