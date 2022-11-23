.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
build: 
	# COPY HLS
	cp -r ../../cmd/hls middle_hls 
	cp -r ../../cmd/hls recv_hls
	cp -r ../../cmd/hls send_hls
	# MAKEFILE BUILD
	make build -C middle_hls
	make build -C recv_hls
	make build -C send_hls
run: 
	# MAKEFILE RUN
	make run -C middle_hls
	make run -C recv_hls
	make run -C send_hls
clean:
	# REMOVE HLS 
	rm -rf middle_hls/hls recv_hls/hls send_hls/hls
	# MAKEFILE CLEAN 
	make clean -C middle_hls
	make clean -C recv_hls
	make clean -C send_hls
