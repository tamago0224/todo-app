FROM golang:1.19.3

WORKDIR /app

COPY . /app

RUN go mod tidy && go build .  && \
    apt-get update && apt-get install -y sqlite3 && \
    sqlite3 todo.db < sql/01-initial-schema.sql

ENTRYPOINT [ "/app/golang-rest-api" ]