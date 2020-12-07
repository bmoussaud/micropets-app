FROM golang
RUN go get -d -v -u github.com/magiconair/properties
EXPOSE 7007
ADD . /go/src/xebialabs.demo/micropetportal/fishes
RUN go install xebialabs.demo/micropetportal/fishes
ENTRYPOINT /go/bin/fishes
