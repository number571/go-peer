GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt0 &
clean:
	pkill -15 prog_hlt0 || true
	rm -rf prog_hlt0 hlt.db 
