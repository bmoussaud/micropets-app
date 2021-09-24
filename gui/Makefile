BINARY_NAME=gui
REGISTRY=localregistry:5000

BINARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev

MICROPETS_NS=micropets-supplychain

all: docker-build docker-push

build:
	ng build --base-href /gui/

test:
	ng serve
	
deps:
	npm install --save-dev @angular-devkit/build-angular

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

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/gui/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml

cnb-image:
	pack build $(BINARY_NAME) --env BP_NODE_VERSION=14 --buildpack gcr.io/paketo-buildpacks/nodejs --builder paketobuildpacks/builder:base


configure-kpack:
	ytt --ignore-unknown-comments -v image_prefix=harbor.mytanzu.xyz/library/micropet -f kpack.yml | kapp deploy --yes --into-ns $(MICROPETS_NS) -a micropet-pets-kpack -f-