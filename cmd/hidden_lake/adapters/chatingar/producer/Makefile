GC=go build
BINPATH=../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/chatingar_producer ./cmd/chatingar_producer
run:
	$(BINPATH)/chatingar_producer
clean:
	rm -f hla_chatingar_producer.yml $(BINPATH)/chatingar_producer
