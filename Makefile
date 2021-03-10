
docker-build:
	$(MAKE) -C cats docker-build
	$(MAKE) -C dogs docker-build
	$(MAKE) -C fishes docker-build
	$(MAKE) -C pets-spring docker-build
	$(MAKE) -C gui docker-build