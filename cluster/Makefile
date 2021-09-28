MICROPETS_SP_NS="micropets-supplychain"

namespace:
	kubectl create namespace $(MICROPETS_SP_NS) --dry-run=client -o yaml | kubectl apply -f -
	kubectl get namespace $(MICROPETS_SP_NS) 

kpack: namespace	
	ytt --ignore-unknown-comments --data-values-env  MICROPETS   -f . | kapp deploy --yes --into-ns $(MICROPETS_SP_NS) -a micropet-kpack -f-

tap: namespace
	ytt --ignore-unknown-comments -f tap/app-operator -f tap/developer -f values.yaml | kapp deploy --yes --into-ns $(MICROPETS_SP_NS) -a micropet-tap -f-

watch_workload:
	watch kubectl tree Workload dogs -n  $(MICROPETS_SP_NS)

untap:
	kapp delete --yes -a micropet-tap

