CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/service
	make build -C $(CMDPATH)/hidden_lake/helpers/traffic
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/send
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/recv
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/service
	make build -C ./recv_hls
	# COPY HLS PROGRAM
	cp -r $(BINPATH)/hls ./recv_hls/prog_hls1
	cp -r $(BINPATH)/hls ./send_hls/prog_hls2
	# COPY HLT PROGRAM
	cp -r $(BINPATH)/hlt ./recv_hls/prog_hlt1
	cp -r $(BINPATH)/hlt ./send_hls/prog_hlt2
	# COPY SERVICE 
	cp $(BINPATH)/common_service ./service/prog_service
	# COPY ADAPTERS 
	cp $(BINPATH)/common_send ./recv_hls/prog_send1
	cp $(BINPATH)/common_send ./send_hls/prog_send2
	cp $(BINPATH)/common_recv ./recv_hls/prog_recv1
	cp $(BINPATH)/common_recv ./send_hls/prog_recv2
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
