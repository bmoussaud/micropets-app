REGISTRY=registry.local:5000
BINARY_NAME=pets

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGET_FLAGS=-u -v

INARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev


all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go

test:
	./mvnw test
clean:
	./mvnw clean

run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	GODEBUG=netdns=go+1 SERVICE_CONFIG=pets_config.dev.json ./pets

deps:	
	$(GOGET) $(GOGET_FLAGS) github.com/spf13/viper
	$(GOGET) $(GOGET_FLAGS) github.com/imroc/req


docker-build:
	./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=registry.local:5000/micropet/pets-spring
	docker push registry.local:5000/micropet/pets-spring:latest

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/pets/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml

logs:
	kubectl logs -f $(kubectl get pods --selector "app=pets-spring-app" --output=name)