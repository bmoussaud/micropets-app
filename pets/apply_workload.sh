NS=micropets-dev
kubectl apply -f config/pets_config.yaml --namespace $NS
tanzu apps workload apply -f config/workload.yaml --live-update --local-path . --source-image akseutap4registry.azurecr.io/pets --namespace $NS --yes  --update-strategy merge
