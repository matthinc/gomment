FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y golang-go python3

RUN apt-get install -y ca-certificates make

COPY . /app
WORKDIR /app
RUN go build

RUN apt-get install -y python3-requests

