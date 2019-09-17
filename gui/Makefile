BINARY_NAME=gui
REGISTRY=localhost:5000
DOCKER_IMAGE=localhost:5000/$(BINARY_NAME):$(XL_VALUE_version)

build:
	ng build

test:
	npm serve
clean:	
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	./$(BINARY_NAME)

docker-build:
	eval $(minikube docker-env)
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)

xebialabs:
	xl --xl-deploy-url http://localhost:4525 apply -s --proceed-when-dirty -f xebialabs.yaml

deploy:
	../xld.sh Applications/Containers/micropets/gui/$(XL_VALUE_version) Environments/MicroPet/Ops/micropet.k8s.staging

