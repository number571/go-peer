.PHONY: default build clean

default: build 
build:
	make -C ./service+messenger build
	make -C ./service+traffic build
	make -C ./service+traffic+messenger build
clean:
	make -C ./service+messenger clean
	make -C ./service+traffic clean
	make -C ./service+traffic+messenger clean
