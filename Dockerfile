FROM golang:alpine AS builder
WORKDIR /src

# Install Node.js for webapp build
RUN apk add --no-cache nodejs npm

# Install webapp dependencies (leveraging Docker layer cache)
COPY webapp/package.json webapp/package-lock.json* webapp/
RUN cd webapp && npm install

# Copy source and build webapp (required for go:embed)
COPY . .
RUN cd webapp && npm run build

# Build Go binary (go:embed picks up webapp/build)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/featmap/featmap .

# Minimal runtime image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/featmap/featmap /opt/featmap/featmap

EXPOSE 5000
WORKDIR /opt/featmap
ENTRYPOINT ["./featmap"]
