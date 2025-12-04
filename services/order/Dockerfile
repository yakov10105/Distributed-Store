FROM golang:1.25.5-alpine AS builder
WORKDIR /app

# Copy workspace files
COPY go.work .
COPY pkg/go.mod pkg/go.mod
COPY services/analytics/go.mod services/analytics/go.mod
COPY services/auth/go.mod services/auth/go.mod
COPY services/bff/go.mod services/bff/go.mod
COPY services/notification/go.mod services/notification/go.mod
COPY services/order/go.mod services/order/go.mod
COPY services/shipping/go.mod services/shipping/go.mod

# Copy source code
COPY pkg pkg
COPY services/order services/order

# Build
RUN go build -o main ./services/order

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 50052
CMD ["./main"]
