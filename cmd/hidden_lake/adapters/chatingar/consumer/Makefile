GC=go build
BINPATH=../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/chatingar_consumer ./cmd/chatingar_consumer
run:
	./$(BINPATH)/chatingar_consumer
clean:
	rm -f hla_chatingar_consumer.yml $(BINPATH)/chatingar_consumer
	rm -rf hla_chatingar_consumer.db
