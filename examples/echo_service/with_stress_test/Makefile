CMDPATH=../../../cmd
BINPATH=../../../bin

NUM_NODES=$(shell ls ./other_nodes | wc -l)

.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version

generate-configs: 
	go run ./_generate_configs

build: generate-configs
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/service
	make build -C ./recv_hls
	# COPY HLS
	cp -r $(BINPATH)/hls ./recv_hls/prog_hls1
	cp -r $(BINPATH)/hls ./send_hls/prog_hls2
	for i in $$(seq 1 $(NUM_NODES)); do \
		cp -r $(BINPATH)/hls ./other_nodes/other_hls$${i}/prog_other_hls$${i}; \
	done;

run:
	# MAKEFILE RUN
	make run -C ./recv_hls
	make run -C ./send_hls
	for i in $$(seq 1 $(NUM_NODES)); do \
		make run -C ./other_nodes/other_hls$${i}; \
	done;

clean:
	# MAKEFILE CLEAN 
	make clean -C ./recv_hls
	make clean -C ./send_hls
	for i in $$(seq 1 $(NUM_NODES)); do \
		make clean -C ./other_nodes/other_hls$${i}; \
	done;
	rm -rf other_nodes/
