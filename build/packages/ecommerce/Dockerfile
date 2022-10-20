# syntax=docker/dockerfile:1.4
FROM golang:1.18-buster as dependencies

ARG BUILD_VERSION

RUN go env -w GOPRIVATE=dev.azure.com/martsoftfilmz/*

WORKDIR /app

COPY  ./ ./
COPY  go.mod ./
COPY  go.sum ./

FROM dependencies as test
RUN go install gotest.tools/gotestsum@latest
RUN gotestsum --junitfile testreport-functionaltests.xml ./tests/sociable/ecommerce -v

FROM scratch as export-test-results
COPY --from=test  /app/*.xml .

FROM dependencies AS builder
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
COPY --from=builder /app/config/config.json .
COPY --from=builder /app/config/config.local.json .

EXPOSE 80
ENV GODEBUG=madvdontneed=1
ENV GOGC=50
CMD ["./ecommerce", "-config", "config.local.json"]


