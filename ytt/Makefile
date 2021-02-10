#should be fetched from package.yaml
ENV=demo
PROFILE=profiles/simple
NAMESPACE=mydemo
APPLICATION=micro-pet
VERSION=1.0.0

katapult:
	ytt -f ${PROFILE}  -f ${ENV}

apply:
	ytt -f ${PROFILE}  -f ${ENV} | kubectl apply -n ${NAMESPACE} -f-
	kubectl get pods -l version=${VERSION},application=${APPLICATION} -n ${NAMESPACE}

kapply:	
	ytt -f ${PROFILE}  -f ${ENV} | kapp -y deploy -a ${APPLICATION} -f -
	kapp inspect  -a ${APPLICATION} --tree

delete:
	ytt -f ${PROFILE}  -f ${ENV} | kubectl delete -f-	

kdelete:
	kapp delete -a ${APPLICATION}