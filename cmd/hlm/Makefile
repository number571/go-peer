GC=go build
.PHONY: default build run clean
default: build run
build:
	$(GC) -o hlm *.go
run:
	./hlm
clean:
	rm -rf hlm hlm.cfg
