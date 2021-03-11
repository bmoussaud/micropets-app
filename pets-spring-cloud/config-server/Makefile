REGISTRY=localregistry:5000
BINARY_NAME=config-server


$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-SNAPSHOT
DOCKER_IMAGE=$(REGISTRY)/pets/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_LATEST=$(REGISTRY)/pets/$(BINARY_NAME):latest
all: test build

build:
	gradle build

test:
	gradle test

clean:
	gradle clean

docker-build:
	gradle bootBuildImage
	docker tag  $(BINARY_NAME):$(VERSION) $(DOCKER_IMAGE_LATEST)
	docker push $(DOCKER_IMAGE_LATEST)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

docker-run: 
	docker run -v config-repo:/config-repo -p 8888:8888 -e SPRING_SECURITY_USER_NAME=scott -e SPRING_SECURITY_USER_PASSWORD=tiger -e SPRING_PROFILES_ACTIVE=docker,git $(BINARY_NAME):$(VERSION)

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl create -f k8s/resources-dev.yaml
	kubectl get pods

logs:
	kubectl logs -f $(kubectl get pods --selector "app=pets-spring-app" --output=name)
