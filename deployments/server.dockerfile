FROM golang:1.24-alpine AS builder

ARG APP_VERSION="dev"
ARG GIT_COMMIT_SHA="unkown"

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -a -tags netgo -ldflags="-w -s -extldflags '-static' -X soccer-api/internal/config.applicationVersion=${APP_VERSION} -X soccer-api/internal/config.gitCommit=${GIT_COMMIT_SHA}" \
    -o /app/server ./cmd/app/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /app/server
COPY --from=builder /app/configs/.env /app/configs/.env

USER 1001:1001

ENTRYPOINT ["/app/server"]