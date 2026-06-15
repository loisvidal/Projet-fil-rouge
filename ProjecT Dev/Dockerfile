FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o RedProject .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/RedProject .
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./RedProject"]
