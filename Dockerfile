FROM golang:latest AS bot_builder
LABEL authors="Cameron Honis, Credits to the stockfish developers"


WORKDIR /
RUN git clone https://github.com/official-stockfish/Stockfish.git

WORKDIR /Stockfish/src
RUN make -j profile-build

FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod download

ENV ENV=prod
RUN go build -o main .

WORKDIR /bots
COPY --from=bot_builder /Stockfish/src/stockfish .
ENV STOCKFISH_PATH=/bots/stockfish

WORKDIR /app