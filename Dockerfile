# Stage 1: Build Go binary
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server main.go

# Stage 2: Final image with minimal dependencies
FROM alpine:latest

# Install rsvg-convert and dependencies
RUN apk update && \
    apk add --no-cache librsvg ca-certificates && \
    apk add libc6-compat

WORKDIR /app

COPY --from=builder /app/server .
COPY index.html .

EXPOSE 8080

CMD ["./server"]
