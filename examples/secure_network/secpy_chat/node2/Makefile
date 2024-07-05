GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hle_2 -priv=../../_keys/4096/priv_node2.key &
clean:
	pkill -15 prog_hle_2 || true
	rm -rf prog_hle_2 main.py
