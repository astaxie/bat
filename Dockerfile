FROM golang:1.8

ADD . /go/src/github.com/astaxie/bat

RUN go install github.com/astaxie/bat

ENTRYPOINT ["/go/bin/bat"]