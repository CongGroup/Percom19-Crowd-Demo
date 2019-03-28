FROM golang:1.10

WORKDIR /go-project/src/Percome19-Crowd-Demo/back-end

ENV GOPATH /go-project

RUN go get -u github.com/golang/dep/cmd/dep

COPY . .

RUN /go-project/bin/dep ensure