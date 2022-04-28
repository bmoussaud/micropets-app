
BINARY_NAME=birds
REGISTRY=localregistry:5000
MICROPETS_NS=micropets-supplychain
MICROPETS_RUN_NS=micropets-run
IMAGE_PREFIX=harbor.mytanzu.xyz/library/micropet
LOCAL_DOCKER_IMAGE=micropet-$(BINARY_NAME)-local-build

run-database:
	docker run -p 5432:5432 --name bird-postgres -e POSTGRES_BD=birds -e POSTGRES_PASSWORD=mysecretpassword -d postgres

run-app:
	./mvnw spring-boot:run

cnb-image:
	pack build $(LOCAL_DOCKER_IMAGE) --env BP_JVM_VERSION=17 --buildpack gcr.io/paketo-buildpacks/java --builder paketobuildpacks/builder:base

workload:
	kubectl apply -f workload.yaml -n $(MICROPETS_NS)

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