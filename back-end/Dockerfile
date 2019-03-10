FROM golang:1.10

WORKDIR /app

ENV GOPATH /app

RUN go get -u "github.com/gorilla/handlers";\
go get -u "github.com/gorilla/mux";\
go get -u  "github.com/gorilla/websocket";\
go get -u  "github.com/ethereum/go-ethereum" ;\
go get -u "github.com/btcsuite/btcd/btcec";

COPY . .