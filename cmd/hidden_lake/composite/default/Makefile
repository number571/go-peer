.PHONY: default build clean

default: build 
build:
	make -C ./service_messenger build
	make -C ./service_traffic build
	make -C ./service_traffic_messenger build
clean:
	make -C ./service_messenger clean
	make -C ./service_traffic clean
	make -C ./service_traffic_messenger clean
