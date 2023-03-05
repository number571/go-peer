# running node3 with input control
# node1, node2 read messages from node3
# middle used for transport messages among nodes
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
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
	cp ../../../../cmd/hidden_lake/messenger/hlm ./node1_hlm/prog_hlm1
	cp ../../../../cmd/hidden_lake/messenger/hlm ./node2_hlm/prog_hlm2
	# COPY HLS PROGRAM
	cp ../../../../cmd/hidden_lake/service/hls ./node1_hlm/prog_hls1
	cp ../../../../cmd/hidden_lake/service/hls ./node2_hlm/prog_hls2
	# COPY HLT PROGRAM
	cp ../../../../cmd/hidden_lake/traffic/hlt ./node1_hlm/prog_hlt1
	cp ../../../../cmd/hidden_lake/traffic/hlt ./node2_hlm/prog_hlt2
	# COPY ADAPTERS 
	cp ../../../../cmd/hidden_lake/adapters/common/service/service ./service/prog_service
	cp ../../../../cmd/hidden_lake/adapters/common/send/send ./node1_hlm/prog_send1
	cp ../../../../cmd/hidden_lake/adapters/common/send/send ./node2_hlm/prog_send2
	cp ../../../../cmd/hidden_lake/adapters/common/recv/recv ./node1_hlm/prog_recv1
	cp ../../../../cmd/hidden_lake/adapters/common/recv/recv ./node2_hlm/prog_recv2
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
