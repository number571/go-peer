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
	make build -C $(CMDPATH)/hidden_lake/composite
	# MAKEFILE BUILD ADAPTERS
	make build -C $(CMDPATH)/hidden_lake/adapters/common/service 
	make build -C $(CMDPATH)/hidden_lake/adapters/common/producer 
	make build -C $(CMDPATH)/hidden_lake/adapters/common/consumer 
	# COPY HLC PROGRAM
	cp $(BINPATH)/hlc ./node1_hlf/prog_hlc1
	cp $(BINPATH)/hlc ./node2_hlf/prog_hlc2
	# COPY STORAGE 
	mkdir ./node2_hlf/hlf.stg
	cp $(STGPATH)/* ./node2_hlf/hlf.stg
	# COPY ADAPTERS 
	cp $(BINPATH)/common_service ./service/prog_service
	cp $(BINPATH)/common_producer ./node1_hlf/prog_producer1
	cp $(BINPATH)/common_producer ./node2_hlf/prog_producer2
	cp $(BINPATH)/common_consumer ./node1_hlf/prog_consumer1
	cp $(BINPATH)/common_consumer ./node2_hlf/prog_consumer2
	# BUILD HLS, HLF CONFIGS
	make build -C node1_hlf
	make build -C node2_hlf
run: 
	# MAKEFILE RUN
	make run -C service
	make run -C node1_hlf
	make run -C node2_hlf
clean:
	# MAKEFILE CLEAN 
	make clean -C service
	make clean -C node1_hlf
	make clean -C node2_hlf
