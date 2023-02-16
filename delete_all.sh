NS=$1
kubectl delete -f cats/config -n $NS
kubectl delete -f dogs/config -n $NS
kubectl delete -f fishes/config -n $NS
kubectl delete -f pets/config -n $NS
kubectl delete -f gui/config -n $NS
