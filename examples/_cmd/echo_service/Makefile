.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
build: 
	# MAKEFILE BUILD
	make build -C ../../../cmd/hidden_lake/service
	make build -C ./recv_hls
	# COPY HLS
	cp -r ../../../cmd/hidden_lake/service/hls ./middle_hls/prog_hls
	cp -r ../../../cmd/hidden_lake/service/hls ./recv_hls/prog_hls1
	cp -r ../../../cmd/hidden_lake/service/hls ./send_hls/prog_hls2
run: 
	# MAKEFILE RUN
	make run -C ./middle_hls
	make run -C ./recv_hls
	make run -C ./send_hls
clean:
	# MAKEFILE CLEAN 
	make clean -C ./middle_hls
	make clean -C ./recv_hls
	make clean -C ./send_hls
