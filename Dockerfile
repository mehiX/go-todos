#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./main.go

#final stage
FROM alpine:latest
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]
CMD ["--http", "0.0.0.0:80", "--dsn", "test:test@tcp(db)/test"]

LABEL Name=todos Version=0.0.1
EXPOSE 80
