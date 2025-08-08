# Build stage
FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /server /server

RUN chmod +x ./server

CMD ["/server"]