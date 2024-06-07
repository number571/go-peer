STGPATH=../_storage
.PHONY: default run clean stop
default: run 
run:
	cp node1_hlf/hls_copy.yml node1_hlf/hls.yml 
	cp node2_hlf/hls_copy.yml node2_hlf/hls.yml 
	mkdir ./node2_hlf/hlf.stg | true
	cp $(STGPATH)/* ./node2_hlf/hlf.stg
	docker-compose build
	docker-compose up
clean:
	rm -f node1_hlf/hls.yml node2_hlf/hls.yml
	rm -rf ./node1_hlf/hlf.stg ./node2_hlf/hlf.stg
stop:
	docker-compose down 
	docker-compose rm -f
