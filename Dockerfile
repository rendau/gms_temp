FROM golang:1.16 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make



FROM alpine:latest

RUN apk add --no-cache --upgrade ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/cmd/build/. ./

HEALTHCHECK --interval=15s --timeout=10s --retries=2  --start-period=3s CMD curl -f http://localhost/healthcheck || exit 1

CMD ["./svc"]
