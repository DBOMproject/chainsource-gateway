# CI Build
FROM golang:1.19.10-bookworm as builder
COPY ./src /build/src
WORKDIR /build/src
ENV GO111MODULE=on
RUN go test -v -cover ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o Gateway

# Package only executable
FROM alpine:latest
WORKDIR /app/
COPY --from=builder build/src/Gateway .
COPY --from=builder build/src/rest_schema src/rest_schema

ENTRYPOINT ["./Gateway"]