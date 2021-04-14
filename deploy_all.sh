
kustomize build dogs/k8s |  kapp -y deploy -n micropetdev -a dogs -f -
kustomize build cats/k8s |  kapp -y deploy -n micropetdev -a cats -f -
kustomize build fishes/k8s |  kapp -y deploy -n micropetdev -a fishes -f -
kustomize build pets/k8s |  kapp -y deploy -n micropetdev -a pets -f -
kustomize build gui/k8s |  kapp -y deploy -n micropetdev -a gui -f -

#kubectl apply -k dogs/k8s
#kubectl apply -k cats/k8s
#kubectl apply -k fishes/k8s
#kubectl apply -k pets/k8s
#kubectl apply -k gui/k8s
