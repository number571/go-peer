# running node3 with input control
# node1, node2 read messages from node3
# middle used for transport messages among nodes
CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/applications/messenger
	make build -C $(CMDPATH)/hidden_lake/helpers/traffic
	make build -C $(CMDPATH)/hidden_lake/service
	# COPY HLT
	cp -r $(BINPATH)/hlt ./middle_hlt_1/prog_hlt_1
	cp -r $(BINPATH)/hlt ./middle_hlt_2/prog_hlt_2
	cp -r $(BINPATH)/hlt ./middle_hlt_3/prog_hlt_3
	# COPY HLM PROGRAM
	cp $(BINPATH)/hlm ./node2_hlm/prog_hlm2
	cp $(BINPATH)/hlm ./node1_hlm/prog_hlm1
	# COPY HLS PROGRAM
	cp $(BINPATH)/hls ./node2_hlm/prog_hls2
	cp $(BINPATH)/hls ./node1_hlm/prog_hls1
	# BUILD HLS, HLM CONFIGS
	make build -C node1_hlm
	make build -C node2_hlm
run: 
	# MAKEFILE RUN
	make run -C ./middle_hlt_1
	make run -C ./middle_hlt_2
	make run -C ./middle_hlt_3
	make run -C node2_hlm
	make run -C node1_hlm
clean:
	# MAKEFILE CLEAN 
	make clean -C ./middle_hlt_1
	make clean -C ./middle_hlt_2
	make clean -C ./middle_hlt_3
	make clean -C node2_hlm
	make clean -C node1_hlm
