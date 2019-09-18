FROM golang

ADD . /go/src/xebialabs.demo/micropetportal/fishes

RUN go get -d -v -u github.com/magiconair/properties

RUN go install xebialabs.demo/micropetportal/fishes
EXPOSE 7002
ENTRYPOINT /go/bin/fishes
