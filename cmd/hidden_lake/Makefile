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

docker-default: docker-build docker-run
docker-build:
	docker-compose build
docker-run: 
	docker-compose up
docker-clean:
	rm -rf ./_mounted/hls.* ./_mounted/hlt.* ./_mounted/hlm.*
	docker-compose down 
	docker-compose rm -f

composite-default: composite-build 
composite-build: build
	cp _default_cfg/hls.cfg ./service
	cp _default_cfg/hlt.cfg ./traffic
	cp _default_cfg/hlm.cfg ./messenger
composite-run: 
	make run -C ./service & 
	make run -C ./traffic & 
	make run -C ./messenger & 
composite-clean: clean
	pkill -15 hls || true
	pkill -15 hlt || true
	pkill -15 hlm || true
