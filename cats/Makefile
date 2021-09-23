BINARY_NAME=cats
REGISTRY=localregistry:5000

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev
GO111MODULE=auto

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


all: build docker-build docker-push

build: deps
	GO111MODULE=auto $(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	./$(BINARY_NAME)

deps:
	GO111MODULE=auto $(GOGET) github.com/kelseyhightower/envconfig
	GO111MODULE=auto $(GOGET) github.com/magiconair/properties

cnb-image:
	pack build $(BINARY_NAME) --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:base

deploy_kapp:
	ytt --ignore-unknown-comments  -f kapp | kapp deploy --yes --into-ns micropet-test -a micropet-cats-service -f-

undeploy_kapp:
	kapp delete -a micropet-dogs-service 


# Cross compilation
build-windows:
	GO111MODULE=auto CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

docker-build: build
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

docker-run:
	docker run --rm  --name $(BINARY_NAME)  -v $(ROOT_DIR):/config -e SERVICE_CONFIG_DIR=/config -p 7002:7002 $(DOCKER_IMAGE_DEV)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/cats/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -k k8s
	kubectl apply -k k8s

load:
	vegeta attack -targets=targets.txt -name=300qps -rate=300 -duration=15s > results.300qps.bin
	cat results.300qps.bin | vegeta plot > plot.300qps.html