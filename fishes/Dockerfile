FROM golang
RUN GO111MODULE=auto go get -d -v -u github.com/magiconair/properties
EXPOSE 7007
ADD . /go/src/moussaud.org/micropetportal/fishes
RUN GO111MODULE=auto go install moussaud.org/micropetportal/fishes
ENTRYPOINT /go/bin/fishes
