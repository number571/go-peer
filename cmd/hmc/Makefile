GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o hmc *.go
run:
	./hmc
clean:
	rm -rf hmc hmc.cfg hmc.stg hmc.db
