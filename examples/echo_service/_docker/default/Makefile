.PHONY: docker-default docker-run docker-clean
docker-default: docker-run 
docker-run:
	docker-compose build
	docker-compose up
docker-clean:
	rm -rf ./_mounted_middle/hlt.db ./_mounted_recv/hls.db ./_mounted_send/hls.db
	docker-compose down 
	docker-compose rm -f
