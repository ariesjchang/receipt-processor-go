FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o receipt-processor ./cmd/server/main.go

CMD ["/app/receipt-processor"]
