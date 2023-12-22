GC=go build
BINPATH=../../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common_send .
run:
	$(BINPATH)/common_send
clean:
	rm -f $(BINPATH)/common_send
