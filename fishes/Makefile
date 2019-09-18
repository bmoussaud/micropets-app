BINARY_NAME=fishes
REGISTRY=localhost:5000

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_WINDOWS=$(BINARY_NAME).exe
DOCKER_IMAGE=$(REGISTRY)/$(BINARY_NAME):$(XL_VALUE_version)

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
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/kelseyhightower/envconfig
	$(GOGET) github.com/magiconair/properties


# Cross compilation
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

docker-build:
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/fishes/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging
