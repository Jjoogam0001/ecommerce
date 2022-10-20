# FROM golang:alpine
# RUN mkdir /app
# COPY . /app
# WORKDIR /app
# RUN go build -o main . 
# CMD ["/app/main"]

FROM golang:latest as builder

# Build the binary
RUN mkdir /app
COPY . /app
WORKDIR /app
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
COPY --from=builder /app/config/config.json .
COPY --from=builder /app/config/config.local.json .

EXPOSE 80
CMD ["./ecommerce", "-config", "config.local.json"]
