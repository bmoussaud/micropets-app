FROM golang

ADD . /go/src/xebialabs.demo/micropetportal/pets

RUN go get -d -v -u github.com/magiconair/properties
RUN go get -d -v -u github.com/kelseyhightower/envconfig

RUN go install xebialabs.demo/micropetportal/pets
EXPOSE 7004
ENTRYPOINT /go/bin/pets
