FROM golang
EXPOSE 7002
ENTRYPOINT /go/bin/cats

RUN GO111MODULE=auto go get -d -v -u github.com/magiconair/properties
ADD . /go/src/moussaud.org/micropetportal/cats
RUN GO111MODULE=auto go install moussaud.org/micropetportal/cats


