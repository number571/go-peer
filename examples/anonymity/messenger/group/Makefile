ROOPATH=./../../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/anonymity/_keys/4096/priv_node1.key \
		./node1_hlm/hls.key
	cp \
		$(ROOPATH)/examples/anonymity/_keys/4096/priv_node1.key \
		./node2_hlm/hls.key
	cp \
		$(ROOPATH)/examples/anonymity/_keys/4096/priv_node1.key \
		./node3_hlm/hls.key
	cp ./node1_hlm/hls_copy.yml ./node1_hlm/hls.yml 
	cp ./node2_hlm/hls_copy.yml ./node2_hlm/hls.yml 
	cp ./node3_hlm/hls_copy.yml ./node3_hlm/hls.yml 
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hlm/hls.yml ./node2_hlm/hls.yml ./node3_hlm/hls.yml
	rm -f ./node1_hlm/hls.key ./node2_hlm/hls.key ./node3_hlm/hls.key
stop:
	docker-compose down 
	docker-compose rm -f
