.PHONY: default run clean
default: run 
run:
	cp node1_hlm/hls_copy.yml node1_hlm/hls.yml 
	cp node2_hlm/hls_copy.yml node2_hlm/hls.yml 
	cp node3_hlm/hls_copy.yml node3_hlm/hls.yml 
	cp group_hlm/hls_copy.yml group_hlm/hls.yml 
	docker-compose build
	docker-compose up
clean:
	rm node1_hlm/hls.yml node2_hlm/hls.yml node3_hlm/hls.yml group_hlm/hls.yml
	docker-compose down 
	docker-compose rm -f
