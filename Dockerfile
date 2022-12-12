FROM golang:1.19.3 as builder

WORKDIR /app

COPY . /app

RUN go mod tidy && go build -o rest-app . 

FROM golang:1.19.3 as runner

WORKDIR /app
COPY --from=builder /app/rest-app /app

ENTRYPOINT [ "/app/rest-app" ]