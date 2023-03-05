GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_service 6060 &
	sleep 5
clean:
	pkill prog_service || true
	rm -rf prog_service service.db
