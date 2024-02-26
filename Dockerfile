FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod tidy

ENV ENV=prod
RUN go build -o main .