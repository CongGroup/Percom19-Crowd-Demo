FROM golang:1.10

WORKDIR /app

COPY . .

ENV GOPATH=/app

WORKDIR /app/src

EXPOSE 4000

CMD ["go","run","main.go"]