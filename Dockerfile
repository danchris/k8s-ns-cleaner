FROM golang:1.18.6 as builder

# Add a work directory
WORKDIR /app

# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build app
RUN go build -o app

FROM golang:1.18.6 as prod

# Copy built binary from builder
COPY --from=builder app .
ENTRYPOINT ["./app"]
