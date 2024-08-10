GC=go build
.PHONY: default run clean 
default: clean run 
run:
	cp -r ../client_hle/priv.key ./hle.key
	./prog_hle &
clean:
	pkill -15 prog_hle || true
	rm -rf prog_hle
