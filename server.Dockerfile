# stage 1: build the executable
FROM golang:alpine AS builder
WORKDIR /server
COPY server/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# stage 2: deploy the executable
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /server/main .
EXPOSE 8080
CMD ["./main"]
