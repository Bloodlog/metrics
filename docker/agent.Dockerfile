FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOARCH=amd64 go build -o agent ./cmd/agent/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/agent /app/agent

RUN chmod +x /app/agent


CMD ["/app/agent"]