rm cat.exe

#env GOOS=windows GOARCH=386 go get -u github.com/kardianos/service
#env GOOS=windows GOARCH=386 go get -u github.com/magiconair/properties

env GOOS=windows GOARCH=386  go build -o cat.exe -v main.go
xl apply -s --proceed-when-dirty -f xebialabs.yaml

