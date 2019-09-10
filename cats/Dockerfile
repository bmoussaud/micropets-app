FROM golang

ADD . /go/src/xebialabs.demo/micropetportal/cats

RUN go get -d -v -u github.com/magiconair/properties

RUN go install xebialabs.demo/micropetportal/cats
EXPOSE 7002
ENTRYPOINT /go/bin/cats
