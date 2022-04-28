#!/bin/bash
set -x

TANZU_USERNAME=bmoussaud@vmware.com
TANZU_PASSWORD=j[YjNv2T@

helm registry login registry.tanzu.vmware.com --username=${TANZU_USERNAME} --password=${TANZU_PASSWORD}

helm pull oci://registry.tanzu.vmware.com/tanzu-sql-postgres/postgres-operator-chart --version v1.6.0 --untar --untardir /tmp

kubectl create secret docker-registry regsecret \
    --docker-server=https://registry.tanzu.vmware.com/ \
    --docker-username=${TANZU_USERNAME} \
    --docker-password=${TANZU_PASSWORD}

echo "
---
# specify the url for the docker image for the operator, e.g. gcr.io/<my_project>/postgres-operator
operatorImage: registry.tanzu.vmware.com/tanzu-sql-postgres/postgres-operator:v1.6.0

# specify the docker image for postgres instance, e.g. gcr.io/<my_project>/postgres-instance
postgresImage: registry.tanzu.vmware.com/tanzu-sql-postgres/postgres-instance:v1.6.0

# specify the name of the docker-registry secret to allow the cluster to authenticate with the container registry for pulling images
dockerRegistrySecretName: regsecret

# override the default self-signed cert-manager cluster issuer
certManagerClusterIssuerName: postgres-operator-ca-certificate-cluster-issuer

# set the resources for the postgres operator deployment
resources: {}
#  limits:
#    cpu: 100m
#    memory: 128Mi
#  requests:
#    cpu: 100m
#    memory: 128Mi
" > /tmp/oprator.values.yaml


#kubectl create namespace postgres-operator

helm install my-postgres-operator /tmp/postgres-operator/ --namespace=postgres-operator --create-namespace  --values /tmp/oprator.values.yaml --wait 


echo "done"


dbname=$(kubectl get secret birds-db-db-secret -o go-template='{{.data.dbname | base64decode}}')
username=$(kubectl get secret birds-db-db-secret -o go-template='{{.data.username | base64decode}}')
password=$(kubectl get secret birds-db-db-secret -o go-template='{{.data.password | base64decode}}')