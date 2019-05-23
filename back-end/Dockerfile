FROM golang:1.10

WORKDIR /go-project/src/github.com/xxRanger/Percome19-Crowd-Demo/back-end

ENV GOPATH /go-project

COPY vendor vendor/
COPY zcrypto zcrypto/
COPY contract contract/
COPY appClient appClient/
COPY user user/
COPY etc etc/
COPY main.go main.go/