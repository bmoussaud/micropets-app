# https://www.jesuisundev.com/load-testing-en-5-minutes-avec-vegeta/

BINARY_NAME=vegeta_runner
REGISTRY=localregistry:5000
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev

docker-build: 
	echo $(DOCKER_IMAGE)
	docker build -t $(DOCKER_IMAGE) .
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE_DEV)

docker-push: docker-build 
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_IMAGE_DEV)

docker-hub-push: docker-build 
	docker tag $(DOCKER_IMAGE) $(DOCKER_HUB_IMAGE)
	docker tag $(DOCKER_IMAGE) $(DOCKER_HUB_IMAGE_DEV)
	docker push $(DOCKER_HUB_IMAGE)
	docker push $(DOCKER_HUB_IMAGE_DEV)


load:
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=300qps -rate=300 -duration=15s > results.300qps.bin
	cat results.300qps.bin | vegeta plot > plot.300qps.html

burst: 	
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=1500qps -rate=1500 -duration=5s  > results.1500qps.bin
	cat results.1500qps.bin | vegeta plot > plot.1500qps.html

long:
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=1500qps -rate=1500 -duration=600s > results.1500qps.bin
	cat results.1500qps.bin | vegeta plot > plot.1500qps.html
	cat results.1500qps.bin | vegeta report -type=text