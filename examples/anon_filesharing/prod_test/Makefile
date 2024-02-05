# running node3 with input control
# node1, node2 read messages from node3
# middle used for transport messages among nodes
CMDPATH=../../../cmd
BINPATH=../../../bin
STGPATH=../_storage
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/applications/filesharer
	make build -C $(CMDPATH)/hidden_lake/service
	# COPY HLF PROGRAM
	cp $(BINPATH)/hlf ./node2_hlf/prog_hlf2
	cp $(BINPATH)/hlf ./node1_hlf/prog_hlf1
	# COPY HLS PROGRAM
	cp $(BINPATH)/hls ./node2_hlf/prog_hls2
	cp $(BINPATH)/hls ./node1_hlf/prog_hls1
	# COPY STORAGE 
	mkdir ./node2_hlf/hlf.stg
	cp $(STGPATH)/* ./node2_hlf/hlf.stg
	# BUILD HLS, HLF CONFIGS
	make build -C node1_hlf
	make build -C node2_hlf
run: 
	# MAKEFILE RUN
	make run -C node2_hlf
	make run -C node1_hlf
clean:
	# MAKEFILE CLEAN 
	make clean -C node2_hlf
	make clean -C node1_hlf
