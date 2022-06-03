BINARY_NAME=$(current_dir)
REGISTRY=localregistry:5000
MICROPETS_NS=micropets-supplychain
MICROPETS_RUN_NS=micropets-run
IMAGE_PREFIX=harbor.mytanzu.xyz/library/micropet
LOCAL_DOCKER_IMAGE=micropet-$(BINARY_NAME)-local-build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GO111MODULE=auto

#current directory
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


BINARY_WINDOWS=$(BINARY_NAME).exe
$(eval SHA1=$(shell git rev-parse --short HEAD))
VERSION=1.0.0-$(SHA1)
DOCKER_IMAGE=$(REGISTRY)/micropet/$(BINARY_NAME):$(VERSION)
DOCKER_IMAGE_DEV=$(REGISTRY)/micropet/$(BINARY_NAME):dev
DOCKER_HUB_IMAGE=bmoussaud/micropet_$(BINARY_NAME):$(VERSION)
DOCKER_HUB_IMAGE_DEV=bmoussaud/micropet_$(BINARY_NAME):dev



fmt:
	find . -type f -name "*.go" | grep -v "./vendor*" | xargs gofmt -s -w

build: deps
	GO111MODULE=auto $(GOBUILD) -o $(BINARY_NAME) -v service/pets 

test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/pets
	SERVICE_CONFIG_DIR=/config ./$(BINARY_NAME)

deps:
	GO111MODULE=auto $(GOGET) github.com/kelseyhightower/envconfig
	GO111MODULE=auto $(GOGET) github.com/magiconair/properties
	GO111MODULE=auto $(GOGET) k8s.io/client-go/kubernetes
	GO111MODULE=auto $(GOGET) sigs.k8s.io/controller-runtime

# Cross compilation
build-windows:
	GO111MODULE=auto CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v

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

docker-run:
	docker run --rm  --name $(BINARY_NAME)  -v $(ROOT_DIR):/config -e SERVICE_CONFIG_DIR=/config -e MP_OBSERVABILITY.TOKEN=$(TO_TOKEN) -p 9000:9000 $(LOCAL_DOCKER_IMAGE)

cnb-image:
	pack build $(LOCAL_DOCKER_IMAGE) --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:base

deploy-cnb:	
	ytt --ignore-unknown-comments -v image_prefix=$(IMAGE_PREFIX) -f kpack.yml | kapp deploy --yes --into-ns $(MICROPETS_NS) -a micropet-$(BINARY_NAME)-kpack -f-

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

load:
	vegeta attack -targets=targets.txt -name=300qps -rate=300 -duration=15s > results.300qps.bin
	cat results.300qps.bin | vegeta plot > plot.300qps.html

k8s-deploy:
	kubectl delete -k k8s
	kubectl apply -k k8s
