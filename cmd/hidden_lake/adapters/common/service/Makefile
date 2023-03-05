GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o service .
run:
	./service
clean:
	rm -f service
	rm -rf service.db
