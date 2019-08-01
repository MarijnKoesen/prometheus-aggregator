FROM golang:1.11.5 as builder
MAINTAINER Marijn Koesen <github@koesen.nl>
ADD . /build
RUN cd /build && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/prometheus-aggregator



FROM alpine:latest
EXPOSE 8080/udp 9090
COPY --from=builder /bin/prometheus-aggregator /bin/prometheus-aggregator
ENTRYPOINT ["/bin/prometheus-aggregator"]
