.PHONY: default build clean 
default: build 
build:
	make -C ./service all-build
	make -C ./messenger all-build
	make -C ./traffic all-build
	make -C ./adapters build
clean:
	make -C ./service all-clean
	make -C ./messenger all-clean
	make -C ./traffic all-clean
	make -C ./adapters clean
