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
	make build -C $(CMDPATH)/hidden_lake/helpers/traffic
	# COPY HLE
	cp -r $(BINPATH)/hle ./node1/prog_hle_1
	cp -r $(BINPATH)/hle ./node2/prog_hle_2
	# COPY HLT
	cp -r $(BINPATH)/hlt ./service_hlt/prog_hlt
	# COPY CLIENT
	cp -r $(CMDPATH)/secpy_chat/main.py ./node1/main.py
	cp -r $(CMDPATH)/secpy_chat/main.py ./node2/main.py
run: 
	# MAKEFILE RUN
	make run -C ./node1
	make run -C ./node2
	make run -C ./service_hlt
clean:
	# MAKEFILE CLEAN 
	make clean -C ./node1
	make clean -C ./node2
	make clean -C ./service_hlt
