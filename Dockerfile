FROM golang:latest AS bot_builder
LABEL authors="Cameron Honis, Credits to the stockfish developers"


WORKDIR /
RUN git clone https://github.com/official-stockfish/Stockfish.git

WORKDIR /Stockfish/src
RUN make -j profile-build

WORKDIR /
RUN git clone https://github.com/CameronHonis/Mila.git
WORKDIR /Mila
RUN go mod tidy
RUN go mod download
RUN go build -o mila

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
COPY --from=bot_builder /Mila/mila .
ENV MILA_PATH=/bots/mila

WORKDIR /app