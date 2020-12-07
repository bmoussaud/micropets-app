FROM golang
EXPOSE 7003
RUN go get -d -v -u github.com/magiconair/properties
ADD . /go/src/xebialabs.demo/micropetportal/dogs
RUN go install xebialabs.demo/micropetportal/dogs
ENTRYPOINT /go/bin/dogs

