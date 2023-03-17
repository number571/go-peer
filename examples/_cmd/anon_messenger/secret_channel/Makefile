# running node3 with input control
# node1, node2 read messages from node3
# middle used for transport messages among nodes
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C ../../../../cmd/hidden_lake/messenger
	make build -C ../../../../cmd/hidden_lake/service
	make build -C ../../../../cmd/hidden_lake/traffic
	# MAKEFILE BUILD ADAPTERS
	make build -C ../../../../cmd/hidden_lake/adapters/common/service 
	make build -C ../../../../cmd/hidden_lake/adapters/common/send 
	make build -C ../../../../cmd/hidden_lake/adapters/common/recv 
	# COPY HLM PROGRAM
	cp ../../../../bin/hlm ./node1_hlm/prog_hlm1
	cp ../../../../bin/hlm ./node2_hlm/prog_hlm2
	# COPY HLS PROGRAM
	cp ../../../../bin/hls ./node1_hlm/prog_hls1
	cp ../../../../bin/hls ./node2_hlm/prog_hls2
	# COPY HLT PROGRAM
	cp ../../../../bin/hlt ./node1_hlm/prog_hlt1
	cp ../../../../bin/hlt ./node2_hlm/prog_hlt2
	# COPY ADAPTERS 
	cp ../../../../bin/common_service ./service/prog_service
	cp ../../../../bin/common_send ./node1_hlm/prog_send1
	cp ../../../../bin/common_send ./node2_hlm/prog_send2
	cp ../../../../bin/common_recv ./node1_hlm/prog_recv1
	cp ../../../../bin/common_recv ./node2_hlm/prog_recv2
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
