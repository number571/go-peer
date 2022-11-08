# running node3 with input control
# node1, node2 read messages from node3
# middle used for transport messages among nodes
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill --version
build: 
	# COPY HLS
	cp -r ../../cmd/hls middle_hls 
	cp -r ../../cmd/hls node1_hlm 
	cp -r ../../cmd/hls node2_hlm 
	cp -r ../../cmd/hls node3_hlm
	# COPY HLM 
	cp -r ../../cmd/hlm node1_hlm 
	cp -r ../../cmd/hlm node2_hlm 
	cp -r ../../cmd/hlm node3_hlm
	# MAKEFILE BUILD
	make build -C middle_hls
	make build -C node1_hlm
	make build -C node2_hlm
	make build -C node3_hlm
run: 
	# MAKEFILE RUN
	make run -C middle_hls
	make run -C node1_hlm
	make run -C node2_hlm
	make run -C node3_hlm
clean:
	# REMOVE HLS 
	rm -rf middle_hls/hls node1_hlm/hls node2_hlm/hls node3_hlm/hls 
	# REMOVE HLM 
	rm -rf node1_hlm/hlm node2_hlm/hlm node3_hlm/hlm 
	# MAKEFILE CLEAN 
	make clean -C middle_hls
	make clean -C node1_hlm
	make clean -C node2_hlm
	make clean -C node3_hlm
