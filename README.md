# MicroPets

## Overview

MicroPet is a MicroService Application that includes 4 components:

- `Dogs` is a service managing dogs (Go)
- `Cats` is a service managing cats (Go)
- `Fishes` is a service managing fishes (Go)
- `Pets` is a front ends service managing cats,dogs & fishes (Go)
- `Gui` is a frontend of the wonderfull application (Angular)

All the services are built into a Docker Images
All the service have deployed in a Kubernetes Cluster following the pattern:

Ingress <--> Service <--> Deployement <--> {ConfigMap,Secrets}

![Architecture](img/micropets-msa-2.png)

Note: All the procedure has tested only on Mac using

- Docker For Mac
- [K3D](https://k3d.io/) / [K3S](https://k3s.io/)
- Helm
- Traefik

## Setup the infrastructure

### New Docker Registry

Create a new Docker Registry locally using docker using `registry.local` as DNS name.

```bash
$./k3s/new-docker-registry.sh
```

Edit your local hostname config /etc/hosts

```bash
127.0.0.1 registry.local
```

### Test the registry

```bash
docker pull containous/whoami
docker tag  containous/whoami registry.local:5000/containous/whoami:latest
docker push registry.local:5000/containous/whoami:latest
```

### New K3S Cluster

Create new K3S cluster using the docker registry created previously.
It deploys [Helm](https://helm.sh/) & [Traefik](https://doc.traefik.io/traefik/).

Edit `k3s/new-local-cluster.sh` and set the value for

- CLUSTER_NAME
- K3S_HOME

```bash
$k3s/new-local-cluster.sh
```

### Test k3s configuration

Apply the following configuration

```bash
kubectl apply -f k3s/test-k3s-traefik-contif.yaml
```

and check with your browser you can connect to `https://localhost:80/whoami/` or running

```bash
curl -k https://localhost:80/whoami/
```

![Components](img/components.png)

## Deployments

### Deploy the Dev environment

```bash
kubectl apply -f k8s/resources-dev.yaml
kubectl delete -f k8s/resources-dev.yaml
```

#### Deploy

```bash
K8S_NS='default'
kubectl apply -f dogs/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f cats/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f fishes/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f pets/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f gui/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f pets-spring/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f pets-spring/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl apply -f pets-spring-cloud/config-server/k8s/resources-dev.yaml -n ${K8S_NS}
```

Check the output and the status of the resources to be sure everything is ok.
It's implicit advice you have to follow after running any _kubectl_ command.
_kubectl apply_ is an asynchronous command that returns 'OK' 90% of the time because it asks the Kubernetes Cluster to change the resources' state.

open the website

```bash
open http://gui.dev.pet-cluster.demo/
```

#### Modify the frontend

- edit `gui/src/app/pets/pets.component.css` and change one color
- commit your code
- run `cd gui && make docker-build k8s-deploy`

the Makefile handles :

- the build of the Angular application
- the build of the Docker Image,
- the deployment into the Kubernetes Cluster

#### Undeploy

```bash
K8S_NS='default'
kubectl delete -f dogs/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f cats/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f fishes/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f pets/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f gui/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f pets-spring/k8s/resources-dev.yaml -n ${K8S_NS}
kubectl delete -f pets-spring-cloud/config-server/k8s/resources-dev.yaml -n ${K8S_NS}
```

### Switch pets configuration

Switch between 2 services (dogs & cats) and 3 services (dogs, cats & fishes).

It uses the [kustomize](https://kustomize.io/) to generate the config map and to link it to the deployment.

```bash
cd pets
kubectl delete -k ./kustomize/overlays/2
kubectl apply -k ./kustomize/overlays/2
kubectl apply -k ./kustomize/overlays/3
kubectl apply -k ./kustomize/overlays/2
```

### New environment : test

Target an existing namespace (test) and modify the Ingress resources to use `test` in it.

It uses the [kustomize](https://kustomize.io/)

- to generate the config map,
- to manage the namespace name,
- to change the Ingress URL (from ._dev_.pet-cluster.demo to ._test_.pet-cluster.demo)

```bash
kubectl create ns test
kustomize build  kustomize/test | sed "s/DEV/TEST/g" | sed "s/pets.dev.pet-cluster.demo/pets.test.pet-cluster.demo/g" | kubectl apply -f -
open http://gui.test.pet-cluster.demo/
kustomize build  kustomize/test | kubectl delete -f -
```

### Progressive deployment on Prod

Target an existing namespace (prod) and modify the Ingress resources to use `prod` in it.

```bash
kubectl create ns prod
kustomize build  kustomize/prod | sed "s/DEV/PRODUCTION/g" | sed "s/pets.dev.pet-cluster.demo/pets.prod.pet-cluster.demo/g" | kubectl apply -f -
# the pet application returns dogs & cats
open http://gui.prod.pet-cluster.demo
kubectl apply -f kustomize/prodfish/resources-dev-3.yaml -n prod
# the pet application returns dogs & cats & FISHES ==alternatively==
open http://gui.prod.pet-cluster.demo
kubectl delete -f kustomize/prodfish/resources-dev-2.yaml -n prod

# the pet application returns dogs & cats & FISHES
open http://gui.prod.pet-cluster.demo

kustomize build  kustomize/prod | kubectl delete -f -
```

Use a dedicated configuration to have the 2 versions of the pets implementation (one with fish, one without)

## Canary Deployment using Traefik

```bash
kustomize build kustomize/canary | sed "s/DEV/CANARY/g" | sed "s/pets.dev.pet-cluster.demo/pets.canary.pet-cluster.demo/g" | kubectl apply -f -
```

open http://pets.canary.pet-cluster.demo/

Inject traffic using [slow_cooker](https://github.com/BuoyantIO/slow_cooker)

```bash
./slow_cooker_darwin -qps 100 http://pets.canary.pet-cluster.demo/
```

After few minutes apply `kustomize/canary/pets_activate_20_80.yaml` to have V2 service (20%) and V3 service (80%)
After few minutes apply `kustomize/canary/pets_activate_00_100.yaml` to have V3 service (100%)

if linkerd has been installed you can look at the Grafana Dashbord showing 1/5 of the requests to the _pets_ service goes to v3 including the fishes.

```bash
linkerd dashboard &
linkerd -n canary stat deploy
```

![Result](img/canary2.png)

## Gitops with Flux

Ref [https://toolkit.fluxcd.io/](https://toolkit.fluxcd.io/)
export GITHUB_TOKEN=<your-token>
export GITHUB_USER=<your-username>

flux bootstrap github \
 --owner=$GITHUB_USER \
 --repository=fleet-infra \
 --branch=main \
 --path=staging-cluster \
 --personal

this command creates or updates the `fleet-infra` personal private repository with a new path `staging-cluster` generating the flux configuration for this cluster. Then it applies them to the current cluster.

```bash
.
├── README.md
└── staging-cluster
    └── flux-system
        ├── gotk-components.yaml
        ├── gotk-sync.yaml
        └── kustomization.yaml

2 directories, 4 files
```

```bash
❯ kubectl get pods --namespace flux-system
NAME                                       READY   STATUS    RESTARTS   AGE
helm-controller-6765c95b47-nsmbz           1/1     Running   0          89s
notification-controller-694856fd64-jctjg   1/1     Running   0          89s
source-controller-5bdb7bdfc9-djxgs         1/1     Running   0          89s
kustomize-controller-7f5455cd78-jz6wc      1/1     Running   0          89s
```

Add the `micropets-app` project to be managed.

```bash
flux create source git webapp \
  --url=https://github.com/bmoussaud/micropets-app \
  --branch=master \
  --interval=30s \
  --export > ./staging-cluster/micropets-source.yaml
```

Add the `cats` projects

```bash
flux create kustomization cats \
  --source=webapp \
  --path="./cats/k8s" \
  --prune=true \
  --validation=client \
  --interval=1h \
  --export > ./staging-cluster/cats.yaml
```

Commit & push the 2 resources.

```bash
git add -A && git commit -m "add staging webapp" && git push
```

Check the cats services is now deployed on the `dev` namespace

```bash
❯ flux get kustomizations
NAME       	READY	MESSAGE                                                          	REVISION                                       	SUSPENDED
cats       	True 	Applied revision: master/5f3500cc01bb04c743d80c24162221125566168f	master/5f3500cc01bb04c743d80c24162221125566168f	False
flux-system	True 	Applied revision: main/f41fcc19c9a3c011b527e8d58feb78bf642badfd  	main/f41fcc19c9a3c011b527e8d58feb78bf642badfd  	False
```

Replicate the configuration for the following services: dogs,gui,pets

```bash
git add -A && git commit -m "add staging dogs,pets,gui" && git push
watch flux get kustomizations
```

```bash
Every 2,0s: flux get kustomizations                                                                          MacBook-Pro-de-Benoit.local:

NAME            READY   MESSAGE                                                                 REVISION                                        SUSPENDED
cats            True    Applied revision: master/5f3500cc01bb04c743d80c24162221125566168f       master/5f3500cc01bb04c743d80c24162221125566168f False
flux-system     True    Applied revision: main/5108bf3cf14826770c63ad28e13434841b40079e         main/5108bf3cf14826770c63ad28e13434841b40079e   False
dogs            True    Applied revision: master/5f3500cc01bb04c743d80c24162221125566168f       master/5f3500cc01bb04c743d80c24162221125566168f False
gui             True    Applied revision: master/5f3500cc01bb04c743d80c24162221125566168f       master/5f3500cc01bb04c743d80c24162221125566168f False
pets            True    Applied revision: master/5f3500cc01bb04c743d80c24162221125566168f       master/5f3500cc01bb04c743d80c24162221125566168f False
```

to force Flux to reconcile with repos

```bash
flux reconcile kustomization flux-system --with-source
```

The `interval` values in the kustomization resource creation tells fluxcd will check this interval to keep the configuration synchronized.
In the command line we set 1h, but it's possible to lower this value to 60s.

- Edit the cats.yaml file and set the new interval value

```yaml
interval: 60s
```

- commit and push the code
- kill the `cats` deployment object.
- wait for 60 secondes, it will come back.

## Pets Spring

[Pet-Spring](pet-spring) is another implementation of the `pets` service using the Spring Cloud Framework.
IPets has been implemented using the Go Language.
It uses the custom `config-server` stored in `pets-spring-cloud/config-server` the configuration is stored in `https://github.com/bmoussaud/sample_configuration`

## YTT (Katapult)

the [ytt](ytt) folder demoes the usage of [ytt](https://carvel.dev/ytt/) & [kapp](https://carvel.dev/kapp/) from the Carvel.dev project to generate the YAML definition of all the resources and apply them into a kubernetes namespace.

From the ytt folder, run:

- `make katapult` generates the YAML K8S resource
- `make kapp` generates the YAML K8S resource and run kapp.
- `make kdelete` deletes the application

## References

- https://blog.stack-labs.com/code/kustomize-101/
- https://kubectl.docs.kubernetes.io/references/kustomize/
- https://tasdikrahman.me/2019/09/12/ways-to-do-canary-deployments-kubernetes-traefik-istio-linkerd/
- https://medium.com/@trlogic/linkerd-traffic-split-acf6fae3b7b8
- https://youtu.be/R6OeIgb7lUI
- https://github.com/JoeDog/siege https://www.linode.com/docs/guides/load-testing-with-siege/
- https://github.com/PacktPublishing/Hands-On-Microservices-with-Spring-Boot-and-Spring-Cloud/tree/master/Chapter12
