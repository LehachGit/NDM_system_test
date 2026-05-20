# syntax=docker/dockerfile:1

FROM golang:1.22-alpine3.19 AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /xff-app .

FROM alpine:3.19
RUN adduser -S appuser
USER appuser
COPY --from=builder /xff-app /xff-app
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/xff-app"]
