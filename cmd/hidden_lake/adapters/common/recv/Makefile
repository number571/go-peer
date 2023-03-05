GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o recv .
run:
	./recv
clean:
	rm -f recv
	rm -rf recv.db
