FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.23

RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]