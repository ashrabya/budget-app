# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /budget-app .

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /budget-app .

# Create data directory
RUN mkdir -p /data

ENV PORT=8080
ENV DATA_FILE=/data/budget-data.json

EXPOSE 8080

ENTRYPOINT ["/app/budget-app"]
