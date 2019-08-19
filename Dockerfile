FROM golang:latest
WORKDIR /build
COPY vendor /build/vendor/
COPY go.mod go.sum *.go /build/
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=0 /build/cfdyndns /root/
CMD ["./cfdyndns"]
