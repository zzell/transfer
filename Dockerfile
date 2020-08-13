FROM golang:1.13.10-alpine as builder

ENV GO111MODULE=on
WORKDIR /transfer
COPY . .
RUN go mod download
RUN go build -o transfer

FROM alpine:latest
COPY --from=builder /transfer/transfer .
COPY --from=builder /transfer/config.json .
