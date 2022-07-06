FILES = main.go settings.go logger.go

.PHONY: default build clean
default: clean build 
build: $(FILES)
	go build $(FILES)
clean:
	rm -rf main chain\:*
