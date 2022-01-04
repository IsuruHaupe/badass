FROM golang:1.16-alpine
WORKDIR /app

COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

COPY server/ .

RUN go build -o main

# wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead

ENTRYPOINT ["./main"]
