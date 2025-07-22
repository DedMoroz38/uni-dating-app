FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY . .

RUN /go/bin/sqlc generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /main /main

EXPOSE 3000

CMD ["/main"] 