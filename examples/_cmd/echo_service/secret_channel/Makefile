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
	cp -r ../../../../cmd/hidden_lake/service/hls ./recv_hls/prog_hls1
	cp -r ../../../../cmd/hidden_lake/service/hls ./send_hls/prog_hls2
	# COPY HLT PROGRAM
	cp -r ../../../../cmd/hidden_lake/traffic/hlt ./recv_hls/prog_hlt1
	cp -r ../../../../cmd/hidden_lake/traffic/hlt ./send_hls/prog_hlt2
	# COPY SERVICE 
	cp ../../../../cmd/hidden_lake/adapters/common/service/service ./service/prog_service
	# COPY ADAPTERS 
	cp ../../../../cmd/hidden_lake/adapters/common/send/send ./recv_hls/prog_send1
	cp ../../../../cmd/hidden_lake/adapters/common/send/send ./send_hls/prog_send2
	cp ../../../../cmd/hidden_lake/adapters/common/recv/recv ./recv_hls/prog_recv1
	cp ../../../../cmd/hidden_lake/adapters/common/recv/recv ./send_hls/prog_recv2
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
