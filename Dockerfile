FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY cluster-metrics /usr/local/bin/
RUN chmod +x /usr/local/bin/cluster-metrics
ENTRYPOINT ["/usr/local/bin/cluster-metrics"]
CMD ["--help"]
