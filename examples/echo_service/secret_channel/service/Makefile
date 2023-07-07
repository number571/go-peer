GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_service 6060 &
clean:
	pkill -15 prog_service || true
	rm -rf prog_service common_service.db
