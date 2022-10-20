FROM golang:latest as builder

# Build the binary
WORKDIR /app
COPY ./ ./
RUN cp ./go.mod /.
RUN cp ./go.sum /.
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/ecommerce/

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
COPY --from=builder /app/config/config.local.json .

EXPOSE 80
CMD ["./ecommerce", "-config", "config.local.json"]
