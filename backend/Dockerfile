# Build Stage
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/api/main.go

# Run Stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
