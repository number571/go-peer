GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o ../../../../../bin/common_send .
run:
	./../../../../../bin/common_send
clean:
	rm -f ../../../../../bin/common_send
