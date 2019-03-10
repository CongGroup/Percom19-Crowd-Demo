FROM golang:1.10

WORKDIR /app

COPY . .

RUN go get -u "github.com/gorilla/handlers";\
"github.com/gorilla/mux";\
"github.com/gorilla/websocket";\
"github.com/ethereum/go-ethereum;"