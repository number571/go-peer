.PHONY: default build clean

default: build 
build:
	make -C ./default/service+messenger build
	make -C ./default/service+traffic build
	make -C ./default/service+traffic+messenger build
	make -C ./mobile/service+messenger build
clean:
	make -C ./default/service+messenger clean
	make -C ./default/service+traffic clean
	make -C ./default/service+traffic+messenger clean
	make -C ./mobile/service+messenger clean
