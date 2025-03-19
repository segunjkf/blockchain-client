FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o blockchain-client .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/blockchain-client .

EXPOSE 8080

ENV GIN_MODE=release
ENV PORT=8080

CMD ["./blockchain-client"]