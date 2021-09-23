REGISTRY=localregistry:5000
BINARY_NAME=pets

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGET_FLAGS=

INARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev

all: docker-build

build: deps
	$(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run: build	
	GO111MODULE=auto GODEBUG=netdns=go+1 SERVICE_CONFIG=pets_config.dev.json ./pets

deps:	
	 $(GOGET) github.com/spf13/viper	
	 $(GOGET) github.com/imroc/req

cnb-image:
	pack build $(BINARY_NAME) --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:base

deploy_kapp:
	ytt --ignore-unknown-comments  -f kapp | kapp deploy --yes --into-ns micropet-test -a micropet-pets-service -f-

undeploy_kapp:
	kapp delete -a micropet-pets-service 


# Cross compilation
build-windows:
	GO111MODULE=auto CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

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

docker-run:
	docker run $(DOCKER_IMAGE_DEV) 

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/pets/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml