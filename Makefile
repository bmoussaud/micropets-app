service_cats := cats
service_dogs := dogs
service_fishes := fishes
service_pets := pets
service_gui := gui
services := $(service_cats) $(service_dogs) $(service_fishes) $(service_pets) $(service_gui)

deploy-kapp:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d deploy-kapp; \
        done
undeploy-kapp:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d undeploy-kapp; \
        done

cnb-image:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d cnb-image; \
        done

deploy-cnb:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d deploy-cnb; \
        done

undeploy-cnb:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d undeploy-cnb; \
        done

docker-build:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d docker-build; \
        done

docker-push:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d docker-push; \
        done


k8s-deploy:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d k8s-deploy; \
        done

workload:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d workload; \
        done

unworkload:
	for d in $(services); \
    	do                               \
          $(MAKE) --directory=$$d unworkload; \
        done

deploy-front:
	kubectx  aws-front
	kustomize build kustomize/aws/front  | kapp -y deploy  -a micropets -f -
	kapp inspect -a micropets

deploy-back:
	kubectx aws-back-admin@aws-back
	kustomize build kustomize/aws/back	  | kapp -y deploy  -a micropets -f -
	kapp inspect -a micropets

kill-front-services:
	kubectx  aws-front
	kubectl delete svc cats-service -n micropet-test
	kubectl delete svc dogs-service -n micropet-test
	kubectl delete svc fishes-service -n micropet-test
	kubectl delete deployment front-cats-app -n micropet-test
	kubectl delete deployment front-dogs-app -n micropet-test
	kubectl delete deployment front-fishes-app -n micropet-test

	
undeploy-app:	
	kapp -y delete -a micropets
