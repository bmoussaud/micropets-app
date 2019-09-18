FROM golang

RUN go get -d -v -u github.com/magiconair/properties
ADD . /go/src/xebialabs.demo/micropetportal/fishes

RUN go install xebialabs.demo/micropetportal/fishes
EXPOSE 7007
ENTRYPOINT /go/bin/fishes
