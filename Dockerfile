################################################### STAGE 1
FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./main.go

################################################### STAGE 2
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./

EXPOSE 8080

ENTRYPOINT ["./main"]