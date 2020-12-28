FROM golang:1.14 AS builder
ENV GOPROXY https://goproxy.io
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -mod vendor -o /enforce-no-rancher-special

FROM alpine:3.12
COPY --from=builder /enforce-no-rancher-special /enforce-no-rancher-special
CMD ["/enforce-no-rancher-special"]