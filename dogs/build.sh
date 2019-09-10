rm dog.exe
export XL_VALUE_version=$1


go get -u github.com/kardianos/service
go get -u github.com/magiconair/properties
go build -o dogs -v main.go
#xl apply -s --proceed-when-dirty -f xebialabs.yaml
#:w../xld.sh Applications/.NET/services/dogs/$XL_VALUE_version  Environments/MicroPet/Dev/micropet.dev
