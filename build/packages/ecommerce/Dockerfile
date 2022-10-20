FROM scratch

ARG  BUILD_VERSION

RUN go env -w GOPRIVATE=dev.azure.com/ecommerce/*

WORKDIR /app

COPY ./ ./
COPY go.mod ./
COPY go.sum ./

RUN ls

RUN CGO_ENABLED=0 GOOS=linux go build  ./cmd/ecommerce/

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


EXPOSE 80
CMD ["./ecommerce"]
