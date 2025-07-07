FROM golang:1.23-bookworm AS builder

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /backend /backend

EXPOSE 5000

CMD ["/backend"]
