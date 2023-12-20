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
	make build -C $(CMDPATH)/hidden_lake/messenger
	make build -C $(CMDPATH)/hidden_lake/service
	make build -C $(CMDPATH)/hidden_lake/traffic
	# COPY HLM PROGRAM
	cp $(BINPATH)/hlm ./node1_hlm/prog_hlm1
	cp $(BINPATH)/hlm ./node2_hlm/prog_hlm2
	cp $(BINPATH)/hlm ./node3_hlm/prog_hlm3
	# COPY HLS PROGRAM
	cp $(BINPATH)/hls ./node1_hlm/prog_hls1
	cp $(BINPATH)/hls ./node2_hlm/prog_hls2
	cp $(BINPATH)/hls ./node3_hlm/prog_hls3
	cp $(BINPATH)/hls ./group_hls/prog_hls0
	# COPY HLT PROGRAM
	cp $(BINPATH)/hlt ./middle_hlt/prog_hlt0
	# BUILD HLS, HLM CONFIGS
	make build -C group_hls
	make build -C node1_hlm
	make build -C node2_hlm
	make build -C node3_hlm
run: 
	# MAKEFILE RUN
	make run -C middle_hlt
	make run -C group_hls
	make run -C node1_hlm
	make run -C node2_hlm
	make run -C node3_hlm
clean:
	# MAKEFILE CLEAN 
	make clean -C middle_hlt
	make clean -C group_hls
	make clean -C node1_hlm
	make clean -C node2_hlm
	make clean -C node3_hlm
