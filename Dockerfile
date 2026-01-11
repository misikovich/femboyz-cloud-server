# builder
FROM golang:1.25.5-alpine AS builder

# dependencies
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


# source
COPY . .

# build
ENV CGO_ENABLED=1
RUN go build -o build/server .

# run prod
FROM alpine:latest

RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=builder /app/build/server .
COPY .env .
COPY environment/ .
CMD ["./server"]
