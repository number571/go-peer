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
	make build -C $(CMDPATH)/hidden_lake/service
	# COPY HLM PROGRAM
	cp $(BINPATH)/hlm ./node1_hlm/prog_hlm1
	cp $(BINPATH)/hlm ./node2_hlm/prog_hlm2
	# COPY HLS PROGRAM
	cp $(BINPATH)/hls ./node1_hlm/prog_hls1
	cp $(BINPATH)/hls ./node2_hlm/prog_hls2
	# BUILD HLS, HLM CONFIGS
	make build -C node1_hlm
	make build -C node2_hlm
run: 
	# MAKEFILE RUN
	make run -C node1_hlm
	make run -C node2_hlm
clean:
	# MAKEFILE CLEAN 
	make clean -C node1_hlm
	make clean -C node2_hlm
