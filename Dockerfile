FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:latest as builder

RUN mkdir /CORE-API-Gateway
WORKDIR /CORE-API-Gateway
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o CORE-API-Gateway .

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /CORE-API-Gateway/CORE-API-Gateway .
COPY config/config.yaml /config
WORKDIR /www
COPY www .

WORKDIR /

ENTRYPOINT [ "/CORE-API-Gateway", "-c", "/config/config.yaml" ]