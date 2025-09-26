FROM golang:1.24.7-alpine3.22 AS build
WORKDIR /src
COPY . /src
ARG VERSION
RUN CGO_ENABLED=0 \
  GOOS=linux \
  go build \
  -ldflags="-s -w -X 'main.version=$VERSION'" \
  -o nextdns-exporter

FROM alpine:3.22 AS src
COPY --from=build /src/nextdns-exporter .
ENTRYPOINT ["/nextdns-exporter"]
