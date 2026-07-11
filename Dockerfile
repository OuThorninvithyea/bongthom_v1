# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /src/admin/kaifin-api

RUN apk add --no-cache ca-certificates git tzdata

COPY admin/kaifin-api/go.mod admin/kaifin-api/go.sum ./
RUN go mod download

COPY admin/kaifin-api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kaifin-api .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S app -G app \
	&& touch .env \
	&& chown -R app:app /app

COPY --from=builder /out/kaifin-api /app/kaifin-api
COPY --from=builder /src/admin/kaifin-api/pkg/i18n/localize /app/pkg/i18n/localize

ENV API_HOST=0.0.0.0
ENV API_PORT=9000

USER app

EXPOSE 9000

ENTRYPOINT ["/app/kaifin-api"]
