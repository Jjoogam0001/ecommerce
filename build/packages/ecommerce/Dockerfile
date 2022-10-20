FROM golang:1.18-buster as builder

ARG  BUILD_VERSION
RUN go env -w GOPRIVATE=dev.azure.com/martsoft/*

WORKDIR /app

COPY --link ./ ./
COPY --link go.mod ./
COPY --link go.sum ./

RUN ls

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.BuildVersion=$BUILD_VERSION" ./cmd/ecommerce/

# Deploy into alpine
FROM alpine:latest

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk add --no-cache \
    curl \
    wget \
    bash \
    openssh \
    shadow \
    npm \
    git \
    jq

WORKDIR /app

COPY --from=builder /app/ecommerce .
COPY --from=builder /app/config/config.json ./config/config.json
COPY --from=builder /app/config/config.local.json ./config/config.local.json


EXPOSE 3001
CMD ["./ecommerce"]
