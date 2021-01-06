# CI Build
FROM golang:1.15-alpine as builder
COPY src /build/src
COPY config /build/config
WORKDIR /build/src
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o Gateway

# Package only executable
FROM alpine:latest
WORKDIR /app/
COPY --from=builder build/src/Gateway .
COPY --from=builder build/src/rest_schema src/rest_schema
COPY --from=builder build/config/agent-config.yaml config/

ENTRYPOINT ["./Gateway"]
