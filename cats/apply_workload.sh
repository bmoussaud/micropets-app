tanzu apps workload apply -f config/workload.yaml --live-update --local-path . --source-image akseutap4registry.azurecr.io/cats --namespace dev-tap --yes  --update-strategy merge
