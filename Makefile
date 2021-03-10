
docker-build:
	$(MAKE) -C cats docker-build
	$(MAKE) -C dogs docker-build
	$(MAKE) -C fishes docker-build
	$(MAKE) -C pets-spring docker-build
	$(MAKE) -C gui docker-build

docker-push:
	$(MAKE) -C cats docker-push
	$(MAKE) -C dogs docker-push
	$(MAKE) -C fishes docker-push
	$(MAKE) -C pets-spring docker-push
	$(MAKE) -C gui docker-push