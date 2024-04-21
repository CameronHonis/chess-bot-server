#!/bin/bash


cd /bots/Stockfish/src || exit
make -j build
export STOCKFISH_PATH=/bots/Stockfish/src/stockfish

cd /bots/Mila || exit
go build -o mila
export MILA_PATH=/bots/Mila/mila

cd /app || exit
./main
