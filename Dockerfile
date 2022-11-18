FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /merge-env-to-ini

FROM busybox
COPY --from=builder /merge-env-to-ini /merge-env-to-ini
CMD ["/merge-env-to-ini"]