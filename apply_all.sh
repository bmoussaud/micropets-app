NS=$1
kubectl apply -f cats/config -n $NS
kubectl apply -f dogs/config -n $NS
kubectl apply -f fishes/config -n $NS
kubectl apply -f pets/config -n $NS
kubectl apply -f gui/config -n $NS
