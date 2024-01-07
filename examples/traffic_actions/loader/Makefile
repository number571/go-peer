CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/helpers/traffic
	make build -C $(CMDPATH)/hidden_lake/helpers/loader
	# COPY HLT
	cp -r $(BINPATH)/hlt ./service_hlt_1/prog_hlt_1
	cp -r $(BINPATH)/hlt ./service_hlt_2/prog_hlt_2
	# COPY HLL
	cp -r $(BINPATH)/hll ./service_hll/prog_hll
run: 
	# MAKEFILE RUN
	make run -C ./service_hlt_1
	make run -C ./service_hlt_2
	make run -C ./service_hll
clean:
	# MAKEFILE CLEAN 
	make clean -C ./service_hlt_1
	make clean -C ./service_hlt_2
	make clean -C ./service_hll
