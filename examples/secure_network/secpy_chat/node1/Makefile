GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hle_1 -key=../../_keys/4096/priv_node1.key &
clean:
	pkill -15 prog_hle_1 || true
	rm -rf prog_hle_1 main.py
