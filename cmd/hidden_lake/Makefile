.PHONY: default build clean \
	docker-default docker-build docker-run docker-clean

default: build 
build:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
	make -C ./composite build
	make -C ./helpers/loader build
	make -C ./helpers/encryptor build
	make -C ./helpers/adapters build
clean:
	make -C ./service clean
	make -C ./messenger clean
	make -C ./traffic clean
	make -C ./composite clean
	make -C ./helpers/loader clean
	make -C ./helpers/encryptor clean
	make -C ./helpers/adapters clean

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
