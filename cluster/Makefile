MICROPETS_SP_NS="micropets-supplychain"

namespace:
	kubectl create namespace $(MICROPETS_SP_NS) --dry-run=client -o yaml | kubectl apply -f -
	kubectl get namespace $(MICROPETS_SP_NS) 

kpack: namespace	
	ytt --ignore-unknown-comments --data-values-env  MICROPETS   -f . | kapp deploy --yes --into-ns $(MICROPETS_SP_NS) -a micropet-kpack -f-

tap: namespace
	ytt --ignore-unknown-comments -f tap/app-operator -f values.yaml | kapp deploy --yes --dangerous-override-ownership-of-existing-resources --into-ns $(MICROPETS_SP_NS) -a micropet-tap -f-

untap:
	kapp delete --yes -a micropet-tap

local-kapp:
	ytt --ignore-unknown-comments -f tap/kapp -f tap/kapp-values/kapp-local-values.yaml

local-gui-kapp:
	ytt --ignore-unknown-comments -f tap/kapp-gui -f tap/kapp-values/kapp-local-values.yaml

kapp-controler:
	kubectl apply -f kapp-sample.yaml

unkapp-controler:
	kubectl delete -f kapp-sample.yaml