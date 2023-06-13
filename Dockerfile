FROM golang:1.20.3 as builder

WORKDIR /app

COPY . /app

RUN go mod tidy
RUN go build -o rest-app . 

FROM golang:1.20.3 as runner

WORKDIR /app
COPY --from=builder /app/rest-app /app

ENTRYPOINT [ "/app/rest-app" ]
