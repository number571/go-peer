.PHONY: docker-default docker-run docker-clean
docker-default: docker-run 
docker-run:
	cp _mounted_node1/hls_copy.cfg _mounted_node1/hls.cfg 
	cp _mounted_node2/hls_copy.cfg _mounted_node2/hls.cfg 
	docker-compose build
	docker-compose up
docker-clean:
	rm -rf _mounted_node1/hls.db _mounted_node2/hls.db
	rm -rf _mounted_node1/hlm.db _mounted_node2/hlm.db
	rm -rf _mounted_middle/hlt.db
	rm -f _mounted_node1/hls.cfg _mounted_node2/hls.cfg 
	docker-compose down 
	docker-compose rm -f
