# Multi-stage Dockerfile for building and running the portfolio-backend

# Builder stage
FROM golang:1.25-trixie AS builder
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -v -trimpath -ldflags "-s -w" -o /usr/local/bin/portfolio-backend ./

# Runtime stage
FROM scratch AS runtime

COPY --from=builder /usr/local/bin/portfolio-backend /usr/local/bin/portfolio-backend

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/portfolio-backend"]
