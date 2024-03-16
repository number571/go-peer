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
	make build -C $(CMDPATH)/hidden_lake/adapters/common/service
	# COPY HLC PROGRAM
	cp -r $(BINPATH)/hlc ./recv_hls/prog_hlc1
	cp -r $(BINPATH)/hlc ./send_hls/prog_hlc2
	# COPY ADAPTED SERVICE 
	cp $(BINPATH)/common_service ./service/prog_service
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
