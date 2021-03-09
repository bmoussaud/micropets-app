FROM golang
EXPOSE 7003
RUN GO111MODULE=auto go get -d -v -u github.com/magiconair/properties
ADD . /go/src/moussaud.org/micropetportal/dogs
RUN GO111MODULE=auto go install moussaud.org/micropetportal/dogs
ENTRYPOINT /go/bin/dogs

