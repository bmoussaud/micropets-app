ng build
export XL_VALUE_version=$1

xl apply -s --proceed-when-dirty -f xebialabs.yaml
../xld.sh Applications/.NET/services/gui/$XL_VALUE_version  Environments/MicroPet/Dev/micropet.dev

