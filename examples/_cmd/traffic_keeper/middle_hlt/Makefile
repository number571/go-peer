GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hlt &
	sleep 5
clean:
	pkill prog_hlt || true
	rm -rf prog_hlt hlt.db
