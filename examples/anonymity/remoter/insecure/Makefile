ROOPATH=./../../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/anonymity/_keys/4096/priv_node1.key \
		./recv_hlc/hls.key
	cp \
		$(ROOPATH)/examples/anonymity/_keys/4096/priv_node2.key \
		./send_hls/hls.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./recv_hlc/hls.key ./send_hls/hls.key
stop:
	docker-compose down 
	docker-compose rm -f
