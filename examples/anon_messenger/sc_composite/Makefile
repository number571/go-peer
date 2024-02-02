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
	make build -C $(CMDPATH)/hidden_lake/composite
	# MAKEFILE BUILD ADAPTERS
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/service 
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/send 
	make build -C $(CMDPATH)/hidden_lake/helpers/adapters/common/recv 
	# COPY HLC PROGRAM
	cp $(BINPATH)/hlc ./node1_hlm/prog_hlc1
	cp $(BINPATH)/hlc ./node2_hlm/prog_hlc2
	# COPY ADAPTERS 
	cp $(BINPATH)/common_service ./service/prog_service
	cp $(BINPATH)/common_send ./node1_hlm/prog_send1
	cp $(BINPATH)/common_send ./node2_hlm/prog_send2
	cp $(BINPATH)/common_recv ./node1_hlm/prog_recv1
	cp $(BINPATH)/common_recv ./node2_hlm/prog_recv2
run: 
	# MAKEFILE RUN
	make run -C service
	make run -C node1_hlm
	make run -C node2_hlm
clean:
	# MAKEFILE CLEAN 
	make clean -C service
	make clean -C node1_hlm
	make clean -C node2_hlm
