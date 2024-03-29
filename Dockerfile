FROM alpine:latest

RUN apk add --no-cache --upgrade ca-certificates tzdata curl

WORKDIR /app

COPY ./cmd/build/. ./

HEALTHCHECK --start-period=7s --interval=10s --timeout=3s --retries=3 CMD curl -f http://localhost/healthcheck || false

CMD ["./svc"]
