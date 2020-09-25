FROM golang:1.15.2 AS builder

ARG ARCH

COPY . /build

RUN cd /build && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -ldflags="-w -s" -o bmwcd_exporter .

# Real image
FROM scratch

COPY --from=builder /build/bmwcd_exporter /bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 9744

ENTRYPOINT [ "/bin/bmwcd_exporter" ]
