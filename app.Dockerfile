FROM golang:1.16-alpine
WORKDIR /app

COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

COPY server/ .

RUN go build -o main

ENTRYPOINT ["./main"]
