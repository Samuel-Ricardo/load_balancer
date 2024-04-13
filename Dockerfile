FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/farely/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

RUN chmod +x ./main

EXPOSE 8080
CMD ["./main"]
