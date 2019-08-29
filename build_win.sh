rm -rf bin
#env GOOS=windows GOARCH=386 go get -u github.com/kardianos/service
#env GOOS=windows GOARCH=386 go get -u github.com/magiconair/properties

echo "** backend..."
env GOOS=windows GOARCH=386  go build -o bin/backend/backend.exe -v backend/main.go
cp -r backend/config.properties bin/backend

echo "** frontend..."
env GOOS=windows GOARCH=386  go build -o bin/frontend/frontend.exe -v frontend/main.go
cp -r frontend/content bin/frontend/content
cp -r frontend/config.properties bin/frontend
echo "** xl..."
xl apply -s --proceed-when-dirty -f xebialabs.yaml
