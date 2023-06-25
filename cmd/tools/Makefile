.PHONY: default build clean

default: build 
build:
	make -C ./encryptor build
	make -C ./keygen build
	make -C ./pmanager build
	make -C ./storage build
clean:
	make -C ./encryptor clean
	make -C ./keygen clean
	make -C ./pmanager clean
	make -C ./storage clean
