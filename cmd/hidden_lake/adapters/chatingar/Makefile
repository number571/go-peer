GC=go build
BINPATH=../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/chatingar ./cmd/chatingar
run:
	./$(BINPATH)/chatingar
clean:
	rm -f hla_chatingar.yml $(BINPATH)/chatingar
	rm -rf hla_chatingar.db
