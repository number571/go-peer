.PHONY: default run clean stop
default: run 
run:
	cp node1_hlm/hls_copy.yml node1_hlm/hls.yml 
	cp node2_hlm/hls_copy.yml node2_hlm/hls.yml 
	docker-compose build
	docker-compose up
clean:
	rm -f node1_hlm/hls.yml node2_hlm/hls.yml
stop:
	docker-compose down 
	docker-compose rm -f
