N=2
kubectl scale deployment fishes-app --replicas ${N}
kubectl scale deployment cats-app --replicas ${N}
kubectl scale deployment dogs-app --replicas ${N}
