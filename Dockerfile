FROM golang:alpine AS builder
WORKDIR /app

RUN apk add build-base curl xz upx

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

RUN upx app

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app .

CMD ["./app"]
