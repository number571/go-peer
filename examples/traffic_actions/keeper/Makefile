CMDPATH=../../../cmd
BINPATH=../../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/traffic
	# COPY HLT
	cp -r $(BINPATH)/hlt ./service_hlt/prog_hlt
run: 
	# MAKEFILE RUN
	make run -C ./service_hlt
clean:
	# MAKEFILE CLEAN 
	make clean -C ./service_hlt
