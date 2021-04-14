#kubectl delete -k dogs/k8s
#kubectl delete -k cats/k8s
#kubectl delete -k fishes/k8s
#kubectl delete -k pets/k8s
#kubectl delete -k gui/k8s

kapp -y delete -n micropetdev -a dogs
kapp -y delete -n micropetdev -a cats
kapp -y delete -n micropetdev -a fishes
kapp -y delete -n micropetdev -a pets
kapp -y delete -n micropetdev -a gui

