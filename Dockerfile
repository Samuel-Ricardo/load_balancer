FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test ./pkg/config

RUN go build -o main ./cmd/farely/main.go
RUN go build -o demo ./cmd/demo/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/demo .
COPY --from=builder /app/example ./example

RUN chmod +x ./main

ARG CONFIG_PATH
ENV CONFIG_PATH=$CONFIG_PATH 

EXPOSE 8080
#CMD ["./main", "--config-path", $CONFIG_PATH]
#CMD ["tail", "-f", "/dev/null"]
