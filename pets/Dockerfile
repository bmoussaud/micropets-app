FROM golang

ADD . /go/src/xebialabs.demo/micropetportal/pets

RUN go get -d -v -u github.com/spf13/viper

RUN go install xebialabs.demo/micropetportal/pets
EXPOSE 7004
ENTRYPOINT /go/bin/pets
