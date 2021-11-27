# ubuntu:rolling @ 2021-11-27 14:30:00
FROM ubuntu@sha256:cc8f713078bfddfe9ace41e29eb73298f52b2c958ccacd1b376b9378e20906ef

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates make golang-go python3 python3-requests python3-dateutil node-typescript

RUN apt-get -y install git

COPY . /app

WORKDIR /app

RUN go get && go build
