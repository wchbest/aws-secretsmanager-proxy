FROM golang:1.13-alpine AS builder

RUN apk add --update --no-cache ca-certificates

ADD . /go/src/github.com/siticom/aws-secretmanager-proxy

WORKDIR  /go/src/github.com/siticom/aws-secretmanager-proxy
RUN CGO_ENABLED=0 go build -o /go/bin/aws-secretmanager-proxy main.go

FROM scratch

COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=builder /go/bin/aws-secretmanager-proxy /aws-secretmanager-proxy

EXPOSE 8080

ENTRYPOINT ["/aws-secretmanager-proxy"]
