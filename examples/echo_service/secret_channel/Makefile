CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/composite
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/send
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/recv
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/service
	# COPY HLC PROGRAM
	cp -r $(BINPATH)/hlc ./recv_hls/prog_hlc1
	cp -r $(BINPATH)/hlc ./send_hls/prog_hlc2
	# COPY ADAPTERS 
	cp $(BINPATH)/common_service ./service/prog_service
	cp $(BINPATH)/common_send ./recv_hls/prog_send1
	cp $(BINPATH)/common_send ./send_hls/prog_send2
	cp $(BINPATH)/common_recv ./recv_hls/prog_recv1
	cp $(BINPATH)/common_recv ./send_hls/prog_recv2
	# BUILD HLS CONFIGS
	make build -C recv_hls
	make build -C send_hls
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
