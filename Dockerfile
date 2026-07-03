# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /src/admin/admin-api

RUN apk add --no-cache ca-certificates git tzdata

COPY admin/admin-api/go.mod admin/admin-api/go.sum ./
RUN go mod download

COPY admin/admin-api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/admin-api .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S app -G app \
	&& touch .env \
	&& chown -R app:app /app

COPY --from=builder /out/admin-api /app/admin-api
COPY --from=builder /src/admin/admin-api/pkg/i18n/localize /app/pkg/i18n/localize

ENV API_HOST=0.0.0.0
ENV API_PORT=9000

USER app

EXPOSE 9000

ENTRYPOINT ["/app/admin-api"]
