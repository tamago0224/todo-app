FROM golang:1.19.3 as builder

WORKDIR /app

COPY . /app

RUN go mod tidy && go build .  && \
    apt-get update && apt-get install -y sqlite3 && \
    sqlite3 todo.db < sql/01-initial-schema.sql

FROM golang:1.19.3 as runner

WORKDIR /app
COPY --from=builder /app/golang-rest-api /app
COPY --from=builder /app/todo.db /var/rest-app/

ENTRYPOINT [ "/app/golang-rest-api" ]