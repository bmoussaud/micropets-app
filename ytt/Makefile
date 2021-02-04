katapult:
	ytt -f .

apply:
	ytt -f . | kubectl apply -f-
	kubectl get pods  -l version=1.0.0,application=micro-pet

kapply:	
	ytt -f . | kapp -y deploy -a micro-pet -f -
	kapp inspect  -a micro-pet  --tree

delete:
	ytt -f . | kubectl delete -f-	

kdelete:
	kapp delete -a micro-pet