.PHONY: default build clean \
	docker-default docker-build docker-run docker-clean

default: build 
build:
	make -C ./service build
	make -C ./messenger build
	make -C ./traffic build
	make -C ./adapters build
	make -C ./_monolith build
clean:
	make -C ./service clean
	make -C ./messenger clean
	make -C ./traffic clean
	make -C ./adapters clean
	make -C ./_monolith clean

docker-default: docker-build docker-run
docker-build:
	docker-compose build
docker-run: 
	docker-compose up
docker-clean:
	rm -rf ./_mounted/hls.* ./_mounted/hlt.* ./_mounted/hlm.*
	docker-compose down 
	docker-compose rm -f

composite-default: build composite-run
composite-run: 
	make run -C ./service & 
	make run -C ./traffic & 
	make run -C ./messenger & 
composite-clean: clean
	pkill -15 hls || true
	pkill -15 hlt || true
	pkill -15 hlm || true
