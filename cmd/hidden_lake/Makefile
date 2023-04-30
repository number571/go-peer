.PHONY: default build clean 
default: build 
build:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
	make -C ./adapters build
clean:
	make -C ./service all-clean
	make -C ./messenger all-clean
	make -C ./traffic all-clean
	make -C ./adapters clean
