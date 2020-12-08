# README

Note: All the procedure has tested only on Mac using

* Docker For Mac
* K3D
* Helm

## New Docker Registry

Create a new Docker Registry locally  using docker using `registry.local` as DNS name.

```bash
$k3s/new-docker-registry.sh
```

Edit your local hostname config /etc/hosts

```bash
127.0.0.1 registry.local
```

## Test the registry

```bash
docker pull containous/whoami
docker tag  containous/whoami registry.local:5000/containous/whoami:latest
docker push registry.local:5000/containous/whoami:latest
```

## New K3S Cluster

Create new K3S cluster using the docker registry created previously.

Edit `k3s/new-local-cluster.sh` and set the value for

* CLUSTER_NAME
* K3S_HOME

```bash
$k3s/new-local-cluster.sh
```

## Test k3s configuration

Apply the following configuration

```bash
kubectl apply -f k3s/test-k3s-traefik-contif.yaml
```

and check with your browser you can connect to `https://localhost:80/whoami/` or running

```bash
curl -k https://localhost:80/whoami/
```


## Deploy

### kubectl

```bash
kubectl apply -f k8s/resources-dev.yaml
kubectl delete -f k8s/resources-dev.yaml
```

```bash
K8S_NS='default'
kubectl apply -f dogs/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f cats/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f fishes/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f pets/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f gui/k8s/resources-dev.yaml -n ${K8S_NS}
```

```bash
K8S_NS='default'
kubectl delete -f dogs/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f cats/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f fishes/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f pets/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f gui/k8s/resources-dev.yaml -n ${K8S_NS}
```


### kustomize :  switch pets configuration

Switch between 2 services (dogs & cats) and 2 services (dogs, cats & fishes).

```bash
kubectl delete -k ./kustomize/overlays/2
kubectl apply -k ./kustomize/overlays/2
kubectl apply -k ./kustomize/overlays/3
kubectl apply -k ./kustomize/overlays/2
```

### new environment test

Target an existing namespace (test) and modify the url to use `test` in it.

```bash
kustomize build  kustomize/test | kubectl apply -f -
```

## Reference

* https://blog.stack-labs.com/code/kustomize-101/

