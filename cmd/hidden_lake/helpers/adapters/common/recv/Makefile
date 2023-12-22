GC=go build
BINPATH=../../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common_recv .
run:
	./$(BINPATH)/common_recv
clean:
	rm -f $(BINPATH)/common_recv
	rm -rf common_recv.db
