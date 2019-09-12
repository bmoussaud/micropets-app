# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGET_FLAGS=-u -v
BINARY_NAME=pets
BINARY_WINDOWS=$(BINARY_NAME).exe
DOCKER_IMAGE=localhost:5000/$(BINARY_NAME):$(XL_VALUE_version)

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
	$(GOGET) $(GOGET_FLAGS) github.com/spf13/viper


# Cross compilation
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

docker-build:
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml
