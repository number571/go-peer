GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt_2 &
clean:
	pkill -15 prog_hlt_2 || true
	rm -rf prog_hlt_2 hlt.db
