FROM golang:1.22.2-alpine3.18

WORKDIR /app

COPY . .
RUN apk add bash git make musl-dev curl


RUN go mod download

