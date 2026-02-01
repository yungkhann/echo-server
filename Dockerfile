FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY backend/ .

RUN /go/bin/swag init -g cmd/server/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./server"]
