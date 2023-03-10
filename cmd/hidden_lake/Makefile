.PHONY: default build clean 
default: build 
build:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
	make -C ./adapters build
clean:
	make -C ./service clean
	make -C ./messenger clean
	make -C ./traffic clean
	make -C ./adapters clean
