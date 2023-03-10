.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
build: 
	# MAKEFILE BUILD
	make build -C ../../../../cmd/hidden_lake/service
	make build -C ../../../../cmd/hidden_lake/traffic
	make build -C ../../../../cmd/hidden_lake/adapters/common/send
	make build -C ../../../../cmd/hidden_lake/adapters/common/recv
	make build -C ../../../../cmd/hidden_lake/adapters/common/service
	make build -C ./recv_hls
	# COPY HLS PROGRAM
	cp -r ../../../../bin/hls ./recv_hls/prog_hls1
	cp -r ../../../../bin/hls ./send_hls/prog_hls2
	# COPY HLT PROGRAM
	cp -r ../../../../bin/hlt ./recv_hls/prog_hlt1
	cp -r ../../../../bin/hlt ./send_hls/prog_hlt2
	# COPY SERVICE 
	cp ../../../../bin/common_service ./service/prog_service
	# COPY ADAPTERS 
	cp ../../../../bin/common_send ./recv_hls/prog_send1
	cp ../../../../bin/common_send ./send_hls/prog_send2
	cp ../../../../bin/common_recv ./recv_hls/prog_recv1
	cp ../../../../bin/common_recv ./send_hls/prog_recv2
run: 
	# MAKEFILE RUN
	make run -C ./service
	make run -C ./recv_hls
	make run -C ./send_hls
clean:
	# MAKEFILE CLEAN 
	make clean -C ./service
	make clean -C ./recv_hls
	make clean -C ./send_hls
