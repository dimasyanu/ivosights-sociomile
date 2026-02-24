FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ARG UNIQUE=default
COPY . .
RUN echo "Cache bust: $UNIQUE"
RUN CGO_ENABLED=0 go build -o sociomile cmd/app/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/sociomile .
COPY .env.docker .env

CMD ["./sociomile"]