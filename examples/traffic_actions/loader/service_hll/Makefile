GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hll &
clean:
	pkill -15 prog_hll || true
	rm -rf prog_hll
