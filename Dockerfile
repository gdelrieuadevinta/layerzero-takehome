FROM golang:1.21.5 as builder

WORKDIR /app

# Copy the go module files and get the necessary dependencies
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Copy the rest of the application
COPY . ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o price_checker

FROM debian:buster-slim
WORKDIR /root/

# Install ca-certificates and clean up in one layer
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/price_checker .

# Default environment variables for the application
ENV MAIN_CURRENCY=bitcoin
ENV VS_CURRENCY=cny

CMD ["./price_checker"]
