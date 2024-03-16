GC=go build
BINPATH=../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common_consumer ./cmd/common_consumer
run:
	./$(BINPATH)/common_consumer
clean:
	rm -f hla_common_consumer.yml $(BINPATH)/common_consumer
	rm -rf common_consumer.db
