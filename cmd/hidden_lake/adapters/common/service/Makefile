GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o ../../../../../bin/common_service .
run:
	./../../../../../bin/service
clean:
	rm -f ../../../../../bin/common_service
	rm -rf ../../../../../bin/common_service.db
