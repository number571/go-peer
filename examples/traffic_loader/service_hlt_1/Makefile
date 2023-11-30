GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt_1 &
clean:
	pkill -15 prog_hlt_1 || true
	rm -rf prog_hlt_1 hlt.db
