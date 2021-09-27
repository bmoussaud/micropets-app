BINARY_NAME=dogs
REGISTRY=localregistry:5000
MICROPETS_NS=micropets-supplychain

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

all: build docker-build docker-push

build: deps
	GO111MODULE=auto $(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	GO111MODULE=auto $(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	./$(BINARY_NAME)

deps:	
	GO111MODULE=auto $(GOGET) github.com/magiconair/properties

cnb-image:
	pack build dog --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:base

deploy_kapp:
	ytt --ignore-unknown-comments  -f kapp | kapp deploy --yes --into-ns micropet-test -a micropet-dogs-service -f-

undeploy_kapp:
	kapp delete -a micropet-dogs-service 


# Cross compilation
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

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
	docker run --rm  --name $(BINARY_NAME)  -v $(ROOT_DIR):/config -e SERVICE_CONFIG_DIR=/config -p 7003:7003 $(DOCKER_IMAGE_DEV)


xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/dogs/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml


create-harbor-secret:
	kubectl create secret generic harbor-registry-credentials \
    --from-file=.dockerconfigjson=./harbor.json \
    --type=kubernetes.io/dockerconfigjson \
    -n micropet-build