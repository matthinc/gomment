FROM alpine:latest AS builder

RUN apk add go

COPY . /app
WORKDIR /app
RUN go build -o /opt/gomment

# Use alpine for production to keep the image small
FROM alpine:latest
RUN apk --no-cache add ca-certificates bash

# Copy binary
COPY --from=builder /opt/gomment /opt
COPY --from=builder /app/frontend /opt/frontend

# Release gin-gonic build
ENV GIN_MODE=release

# Database directory (sqlite)
RUN mkdir /opt/db
ENV GOMMENT_DB_PATH=/opt/db/gomment.db

WORKDIR /opt
CMD ["./gomment"]
