FROM golang:1.23.3-bullseye AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo ./bin/main.go

FROM ${BASE_RUNTIME_IMAGE}alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main", "serve"]
