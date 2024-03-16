GC=go build
BINPATH=../../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common_producer ./cmd/common_producer
run:
	$(BINPATH)/common_producer
clean:
	rm -f hla_common_producer.yml $(BINPATH)/common_producer
