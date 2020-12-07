FROM golang
EXPOSE 7004
RUN go get -d -v -u github.com/spf13/viper
RUN go get -d -v -u github.com/imroc/req
ADD . /go/src/xebialabs.demo/micropetportal/pets
RUN go install xebialabs.demo/micropetportal/pets
ENTRYPOINT /go/bin/pets
