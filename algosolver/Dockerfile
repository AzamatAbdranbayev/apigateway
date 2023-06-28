FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy
COPY . .
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o algosolver ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/algosolver .
CMD "/app/algosolver"
