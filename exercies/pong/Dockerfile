# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /src

# Copy mod files first for better layer caching
COPY go.mod ./
RUN go mod download

# Copy the rest and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app .

# Run stage
FROM alpine:3.20

# optional: ca certs for HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /out/app /app/app

# common convention
ENV PORT=8080
EXPOSE 8080

CMD ["/app/app"]