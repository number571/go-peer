GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt &
clean:
	pkill -15 prog_hlt || true
	rm -rf prog_hlt hlt.db
