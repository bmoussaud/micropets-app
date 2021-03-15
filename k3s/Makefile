CLUSTER_NAME=pet-cluster

destroy-cluster:
	k3d cluster delete $(CLUSTER_NAME)
	docker network disconnect k3d-$(CLUSTER_NAME) localregistry

create-cluster:
	k3d cluster create $(CLUSTER_NAME) --api-port 127.0.0.1:6443 -p 80:80@loadbalancer -p 443:443@loadbalancer --k3s-server-arg "--no-deploy=traefik" --volume "$(CURDIR)/k3d-registries.yaml:/etc/rancher/k3s/registries.yaml"
	k3d kubeconfig get $(CLUSTER_NAME)
	docker network connect k3d-$(CLUSTER_NAME) localregistry
	kubectl cluster-info

install-traefik:
	helm repo add traefik https://containous.github.io/traefik-helm-chart
	helm install traefik traefik/traefik
	echo "access to the traefix dashboard"
	echo 'kubectl port-forward $(kubectl get pods --selector "app.kubernetes.io/name=traefik" --output=name) 9000:9000'
	echo 'open http://localhost:9000/dashboard/#/'