FROM golang:1.10

WORKDIR /go-project/src/Percome19-Crowd-Demo/client

ENV GOPATH /go-project

RUN mkdir /go-project/bin && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD Gopkg.* main.go ./

RUN  /go-project/bin/dep ensure

COPY . .

