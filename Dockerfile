FROM golang:latest AS build

COPY . /go/src/github.com/skpr/cluster-metrics
WORKDIR /go/src/github.com/skpr/cluster-metrics
RUN make build

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/skpr/cluster-metrics/bin/cluster-metrics /usr/sbin/cluster-metrics
RUN chmod +x /usr/sbin/cluster-metrics
ENTRYPOINT ["/usr/sbin/cluster-metrics"]
CMD ["--help"]
