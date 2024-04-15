GC=go build
BINPATH=../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common ./cmd/common
run:
	./$(BINPATH)/common
clean:
	rm -f hla_common.yml $(BINPATH)/common
	rm -rf hla_common.db
