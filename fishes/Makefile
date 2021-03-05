BINARY_NAME=fishes
REGISTRY=registry.local:5000

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

all: test build

build:
	GO111MODULE=auto $(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	GO111MODULE=auto $(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	GO111MODULE=auto $(GOBUILD) -o $(BINARY_NAME) -v main.go
	./$(BINARY_NAME)

deps:	
	GO111MODULE=auto $(GOGET) github.com/magiconair/properties


# Cross compilation
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

docker-build: build
	echo $(DOCKER_IMAGE)
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE_DEV)
	docker push $(DOCKER_IMAGE_DEV)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/fishes/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml