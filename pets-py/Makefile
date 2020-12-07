# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGET_FLAGS=-u -v
BINARY_NAME=petspy
BINARY_WINDOWS=$(BINARY_NAME).exe
DOCKER_IMAGE=localhost:5000/$(BINARY_NAME):$(XL_VALUE_imageversion)

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	GODEBUG=netdns=go+1 SERVICE_CONFIG=pets_config.dev.json ./pets

deps:	
	pip install --no-cache-dir -r requirements.txt	


# Cross compilation
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

docker-build:
	eval $(minikube docker-env)
	export XL_VALUE_imageversion=`date +%Y%m%d%H%M%S`	
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/pets-py/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging
