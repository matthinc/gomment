FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates make golang-go python3 python3-requests python3-dateutil node-typescript

RUN apt-get -y install git

COPY . /app

WORKDIR /app

RUN go get && go build
