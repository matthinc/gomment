FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates make golang-go python3 python3-requests python3-dateutil node-typescript

COPY . /app
WORKDIR /app
RUN go build
