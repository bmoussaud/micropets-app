rm bin/backend.exe
env GOOS=windows GOARCH=386  go get -u github.com/kardianos/service
env GOOS=windows GOARCH=386  go build -o bin/backend.exe -v backend/main.go
xl apply -f xebialabs.yaml
