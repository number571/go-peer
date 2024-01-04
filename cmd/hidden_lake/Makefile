.PHONY: default build clean \
	docker-default docker-build docker-run docker-clean

default: build 
build:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
	make -C ./composites build
	make -C ./helpers build
clean:
	make -C ./service clean
	make -C ./messenger clean
	make -C ./traffic clean
	make -C ./composites clean
	make -C ./helpers clean

docker-default: docker-build docker-run
docker-build: docker-configs
	docker-compose build
docker-configs:
	cp ./_configs/docker/* ./_mounted
docker-run: 
	docker-compose up
docker-clean:
	rm -rf ./_mounted/*.yml ./_mounted/*.db ./_mounted/*.key
	docker-compose down 
	docker-compose rm -f

build-main:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
