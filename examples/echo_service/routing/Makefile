CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/helpers/traffic
	make build -C $(CMDPATH)/hidden_lake/service
	make build -C ./recv_hls
	# COPY HLS
	cp -r $(BINPATH)/hlt ./middle_hlt_1/prog_hlt_1
	cp -r $(BINPATH)/hlt ./middle_hlt_2/prog_hlt_2
	cp -r $(BINPATH)/hlt ./middle_hlt_3/prog_hlt_3
	cp -r $(BINPATH)/hls ./recv_hls/prog_hls1
	cp -r $(BINPATH)/hls ./send_hls/prog_hls2
run: 
	# MAKEFILE RUN
	make run -C ./middle_hlt_1
	make run -C ./middle_hlt_2
	make run -C ./middle_hlt_3
	make run -C ./recv_hls
	make run -C ./send_hls
clean:
	# MAKEFILE CLEAN 
	make clean -C ./middle_hlt_1
	make clean -C ./middle_hlt_2
	make clean -C ./middle_hlt_3
	make clean -C ./recv_hls
	make clean -C ./send_hls
