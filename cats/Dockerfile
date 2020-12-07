FROM golang
EXPOSE 7002
ENTRYPOINT /go/bin/cats

RUN go get -d -v -u github.com/magiconair/properties
ADD . /go/src/xebialabs.demo/micropetportal/cats
RUN go install xebialabs.demo/micropetportal/cats


