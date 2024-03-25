# Build binary from golang image
FROM golang:alpine AS builder
WORKDIR /server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a Docker multi-stage build to create a lean production image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /server/main .
EXPOSE 8080
CMD ["./main"]
