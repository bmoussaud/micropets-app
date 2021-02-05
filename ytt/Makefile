ENV="demo"

katapult:
	ytt -f profiles  -f ${ENV}

apply:
	ytt -f profiles  -f ${ENV} | kubectl apply -f-
	kubectl get pods -l version=1.0.0,application=micro-pet -n ${ENV}

kapply:	
	ytt -f profiles  -f ${ENV} | kapp -y deploy -a micro-pet -f -
	kapp inspect  -a micro-pet --tree

delete:
	ytt -f profiles  -f ${ENV} | kubectl delete -f-	

kdelete:
	kapp delete -a micro-pet