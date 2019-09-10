FROM golang

ADD . /go/src/xebialabs.demo/micropetportal/dogs

RUN go get -d -v -u github.com/magiconair/properties

RUN go install xebialabs.demo/micropetportal/dogs
EXPOSE 7003
ENTRYPOINT /go/bin/dogs
