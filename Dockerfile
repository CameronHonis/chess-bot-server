FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod download

ENV ENV=prod
RUN go build -o main .

# Download stockfish
WORKDIR /bots
RUN git clone https://github.com/official-stockfish/Stockfish.git

# Download mila
WORKDIR /bots
RUN git clone https://github.com/CameronHonis/Mila.git

WORKDIR /app