.PHONY: default build clean

default: build 
build:
	make -C ./default/service_messenger build
	make -C ./default/service_traffic build
	make -C ./default/service_traffic_messenger build
	make -C ./mobile/service_messenger build
clean:
	make -C ./default/service_messenger clean
	make -C ./default/service_traffic clean
	make -C ./default/service_traffic_messenger clean
	make -C ./mobile/service_messenger clean
