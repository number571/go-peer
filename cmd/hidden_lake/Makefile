.PHONY: default build clean \
	docker-default docker-build docker-run docker-clean

default: build 
build:
	make -C ./service all-build
	make -C ./messenger all-build
	make -C ./traffic all-build
	make -C ./adapters build
clean:
	make -C ./service all-clean
	make -C ./messenger all-clean
	make -C ./traffic all-clean
	make -C ./adapters clean

docker-default: build 
docker-build:
	docker-compose build
docker-run: 
	docker-compose up
docker-clean:
	docker-compose down 
	docker-compose rm -f
