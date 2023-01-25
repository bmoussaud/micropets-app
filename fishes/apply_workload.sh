tanzu apps workload apply -f config/workload.yaml --live-update --local-path . --source-image akseutap4registry.azurecr.io/fishes --namespace dev-tap --yes  --update-strategy merge
