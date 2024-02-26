FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod download

ENV ENV=prod
RUN go build -o main .