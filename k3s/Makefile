CLUSTER_NAME=pet-cluster

destroy-cluster:
	k3d cluster delete $(CLUSTER_NAME)
	docker network disconnect k3d-$(CLUSTER_NAME) localregistry

create-cluster:
	k3d cluster create $(CLUSTER_NAME) --api-port 127.0.0.1:6443 -p 80:80@loadbalancer -p 443:443@loadbalancer --volume "$(CURDIR)/k3d-registries.yaml:/etc/rancher/k3s/registries.yaml"
	k3d kubeconfig get $(CLUSTER_NAME)
	docker network connect k3d-$(CLUSTER_NAME) localregistry
	kubectl cluster-info

install-ingress:
	kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
	kubectl get pods -n projectcontour -o wide