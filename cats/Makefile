# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cats

BINARY_WINDOWS=$(BINARY_NAME).exe
DOCKER_IMAGE=xebialabs.demo/$(BINARY_NAME)

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
