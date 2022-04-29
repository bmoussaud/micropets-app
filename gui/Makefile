BINARY_NAME=$(current_dir)
REGISTRY=localregistry:5000
MICROPETS_NS=micropets-supplychain
MICROPETS_RUN_NS=micropets-run
IMAGE_PREFIX=harbor.mytanzu.xyz/library/micropet
LOCAL_DOCKER_IMAGE=micropet-$(BINARY_NAME)-local-build

BINARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev


#current directory
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

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

k8s-deploy:
	kubectl delete -f k8s/resources-dev.yaml
	kubectl apply -f k8s/resources-dev.yaml

cnb-image: 
# TODO: Test pack build <image> -b paketo-buildpacks/nodejs -b paketo-buildpacks/nginx -e "BP_NODE_RUN_SCRIPTS=build"
	pack build $(LOCAL_DOCKER_IMAGE) --env BP_NODE_VERSION=14 --buildpack gcr.io/paketo-buildpacks/nodejs:0.9.0 --builder paketobuildpacks/builder:base

cnb-image-2:
# WIP TODO: Test pack build <image> -b paketo-buildpacks/nodejs -b paketo-buildpacks/nginx -e "BP_NODE_RUN_SCRIPTS=build"
#	pack build $(LOCAL_DOCKER_IMAGE) --buildpack gcr.io/paketo-buildpacks/nodejs:0.11.1   --env "BP_NODE_RUN_SCRIPTS=build" --env "NODE_ENV=production" --builder paketobuildpacks/builder:base
#	pack build $(LOCAL_DOCKER_IMAGE) --buildpack gcr.io/paketo-buildpacks/nodejs --env "BP_NODE_RUN_SCRIPTS=build" --env "NODE_ENV=development"
	pack build $(LOCAL_DOCKER_IMAGE) --buildpack gcr.io/paketo-buildpacks/nodejs:0.17.1 --env "BP_NODE_RUN_SCRIPTS=build" --env "NODE_ENV=development"  --env "NPM_CONFIG_LOGLEVEL=verbose" --buildpack gcr.io/paketo-buildpacks/nginx:0.6.0  --buildpack paketo-community/staticfile
	
docker-run:
	docker run --rm --interactive --tty --init --env PORT=8080 --publish 8080:8080 --name $(BINARY_NAME) $(LOCAL_DOCKER_IMAGE) 
 
deploy-cnb:	
	ytt --ignore-unknown-comments  -v image_prefix=$(IMAGE_PREFIX) -f kpack.yml | kapp deploy --yes --into-ns $(MICROPETS_NS) -a micropet-$(BINARY_NAME)-kpack -f-

undeploy-cnb:	
	kapp delete -a micropet-$(BINARY_NAME)-kpack 

namespace:
	kubectl create namespace $(MICROPETS_RUN_NS) --dry-run=client -o yaml | kubectl apply -f -
	kubectl get namespace $(MICROPETS_RUN_NS) 

deploy-kapp: namespace
	ytt --ignore-unknown-comments  -f kapp | kapp deploy --yes --into-ns $(MICROPETS_RUN_NS) -a micropet-$(BINARY_NAME)-service -f-

undeploy-kapp:
	kapp delete -y -a micropet-$(BINARY_NAME)-service

workload:
	kubectl apply -f Workload.yaml -n $(MICROPETS_NS)

unworkload:
	kubectl delete Workload $(BINARY_NAME) -n $(MICROPETS_NS)

watch_workload:
	watch kubectl tree Workload $(BINARY_NAME) -n  $(MICROPETS_NS)

desc_workload:
	kubectl describe Workload $(BINARY_NAME) -n  $(MICROPETS_NS)

deliverable:
	kubectl apply -f deliverable.yaml -n $(MICROPETS_NS)

undeliverable:
	kubectl delete Deliverable $(BINARY_NAME) -n $(MICROPETS_NS)