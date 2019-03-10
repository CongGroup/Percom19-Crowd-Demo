FROM golang:1.10

WORKDIR /app

COPY . .

RUN go get -u "github.com/gorilla/handlers";\
go get -u "github.com/gorilla/mux";\
go get -u  "github.com/gorilla/websocket";\
go get -u  "github.com/ethereum/go-ethereum"