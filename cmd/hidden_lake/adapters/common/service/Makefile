GC=go build
BINPATH=../../../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/common_service .
run:
	./$(BINPATH)/service
clean:
	rm -f $(BINPATH)/common_service
	rm -rf $(BINPATH)/common_service.db
