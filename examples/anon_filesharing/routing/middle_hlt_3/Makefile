GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt_3 &
clean:
	pkill -15 prog_hlt_3 || true
	rm -rf prog_hlt_3 hlt.db
