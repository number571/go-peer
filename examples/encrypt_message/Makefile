CMDPATH=../../cmd
BINPATH=../../bin
.PHONY: default init build run clean
default: init clean build run
init: 
	# CHECK PROGRAMS
	pkill -15 --version
build: 
	# MAKEFILE BUILD
	make build -C $(CMDPATH)/hidden_lake/helpers/encryptor
	# COPY HLT
	cp -r $(BINPATH)/hle ./service_hle/prog_hle
run: 
	# MAKEFILE RUN
	make run -C ./service_hle
clean:
	# MAKEFILE CLEAN 
	make clean -C ./service_hle
