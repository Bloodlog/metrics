FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOARCH=amd64 go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN chmod +x /app/server

EXPOSE 8080

CMD ["/app/server"]