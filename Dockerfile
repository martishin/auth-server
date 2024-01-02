# Build stage
FROM golang:1.21.4-alpine3.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux go build -o auth-server cmd/sso/main.go

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/auth-server /app/auth-server
COPY configs /app/configs
EXPOSE 8080
CMD ["/app/auth-server"]
