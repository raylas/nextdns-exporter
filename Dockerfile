FROM golang:1.27.3-alpine3.24 AS build
WORKDIR /src
COPY . /src
ARG VERSION
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -ldflags="-s -w -X 'main.version=$VERSION'" \
    -o nextdns-exporter

FROM alpine:3.24 AS src
COPY --from=build /src/nextdns-exporter .
ENTRYPOINT /nextdns-exporter
