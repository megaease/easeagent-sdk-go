FROM golang:1.19-alpine AS builder

RUN apk --no-cache add make

WORKDIR /opt/meshdemo
COPY . .

RUN make build

# ---

FROM alpine:3.16

WORKDIR /opt/meshdemo

COPY entrypoint.sh /
COPY --from=builder /opt/meshdemo/bin/ /opt/meshdemo/bin/
COPY --from=builder /opt/meshdemo/tls* /opt/meshdemo/

RUN apk --no-cache add tini tzdata && \
    chmod +x /entrypoint.sh /opt/meshdemo/bin/*

ENV PATH /opt/meshdemo/bin:$PATH

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
