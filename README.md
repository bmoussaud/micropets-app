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

## Install & Configure Digital.ai Deploy

* Install a brand new Deploy Server if you don't have one
* Check the smoke test plugin is installed else install it : [https://github.com/xebialabs-community/xld-smoke-test-plugin/releases/download/v1.0.7/xld-smoke-test-plugin-1.0.7.xldp]
* Check the gitops plugin is installed else install it : [https://github.com/xebialabs-community/xld-gitops-plugin/releases/download/v0.1.0-rc.1/xld-gitops-plugin-0.1.0-rc.1.xldp]

* run Deploy 9.8 in a docker container

``` bash
$cd xl
$docker-compose up
$docker network connect k3d-$CLUSTER_NAME xl-deploy
```

* Edit the K3S Cluster configuration:

  * run `k3d kubeconfig get $CLUSTER_NAME` command
  * edit `xebialabs/infrastructure.values` file by replacing

    * `caCert` value with `clusters/cluster/certificate-authority-data` value from the output of `k3d kubeconfig`
    * `apiServerURL` value with `clusters/cluster/server` value from the output of `k3d kubeconfig` (if you use the deploy docker setup, `https://k3d-book-cluster-serverlb:6443` will be the value)
    * `tlsCert` value with `name/user/client-certificate-data` value from the output of `k3d kubeconfig`
    * `tlsPrivateKey` value with `name/user/client-key-data` value from the output of `k3d kubeconfig`

    * or in the xebialabs folder, run `python configure_k3d.py --output infrastructure.xlvals --cluster $CLUSTER_NAME`

* Import all the ci definitions (application environment infrastructure) : run `make initialci`
* Deploy the application using the UI or the command line

```bash
$./xlw preview --values version=1.0.0-0.0.1 -f xebialabs/deployment.yaml
$./xlw apply --values version=1.0.0-0.0.1  -s -p -f xebialabs/deployment.yaml
````

## Modify the application

### Edit the books.html

Edit `src/main/webapp/books.html`  and modify the css -> footer -> background-color (line #24)

### Build the web application

```bash
$make web
````

this command will

* build the war artifact,
* put it into a Docker container `docker build`,
* tag the images `docker tag bmoussaud/bookstore-advanced registry.local:5000/bmoussaud/bookstore-advanced`
* push it into the docker registry `docker push registry.local:5000/bmoussaud/bookstore-advanced`

### Build the database

```bash
$make database
```

or in detail:

```bash
cd database
export SHA1="0.0.1"
export DB_VERSION="1.0.0-$SHA1"
docker build . --tag bmoussaud/bookstore-advanced-database:$DB_VERSION --build-arg version=$DB_VERSION
docker tag bmoussaud/bookstore-advanced-database:$DB_VERSION registry.local:5000/bmoussaud/bookstore-advanced-database:$DB_VERSION
docker push registry.local:5000/bmoussaud/bookstore-advanced-database:$DB_VERSION
```

### Deploy


git commit -am "change the color"

```bash
$make deploy
```




## References

* https://codeburst.io/creating-a-local-development-kubernetes-cluster-with-k3s-and-traefik-proxy-7a5033cb1c2d
* https://k33g.gitlab.io/articles/2020-02-27-K3S-05-REGISTRY.html
* https://k3d.io/usage/guides/registries
* https://blog.ruanbekker.com/blog/2020/02/21/persistent-volumes-with-k3d-kubernetes/
