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
	make build -C ../../../cmd/hls
	make build -C ../../../cmd/hlm
	# COPY HLS PROGRAM
	cp ../../../cmd/hls/hls ./middle_hls/prog_hls
	cp ../../../cmd/hls/hls ./node1_hlm/prog_hls1
	cp ../../../cmd/hls/hls ./node2_hlm/prog_hls2
	# COPY HLM PROGRAM
	cp ../../../cmd/hlm/hlm ./node1_hlm/prog_hlm1
	cp ../../../cmd/hlm/hlm ./node2_hlm/prog_hlm2
	# COPY HLM WEB
	cp -r ../../../cmd/hlm/web ./node1_hlm
	cp -r ../../../cmd/hlm/web ./node2_hlm
run: 
	# MAKEFILE RUN
	make run -C middle_hls
	make run -C node1_hlm
	make run -C node2_hlm
	sleep 5
clean:
	# MAKEFILE CLEAN 
	make clean -C middle_hls
	make clean -C node1_hlm
	make clean -C node2_hlm
