FROM golang
EXPOSE 7004
RUN GO111MODULE=auto  go get -d -v -u github.com/spf13/viper
RUN GO111MODULE=auto  go get -d -v -u github.com/imroc/req
ADD . /go/src/moussaud.org/micropetportal/pets
RUN GO111MODULE=auto  go install moussaud.org/micropetportal/pets
ENTRYPOINT /go/bin/pets
