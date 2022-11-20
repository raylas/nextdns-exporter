FROM golang:1.19.3-alpine3.16 AS build
WORKDIR /src
COPY . /src
RUN CGO_ENABLED=0 GOOS=linux go build -o nextdns-exporter

FROM alpine:3.16 AS src
COPY --from=build /src/nextdns-exporter .
ENTRYPOINT /nextdns-exporter
