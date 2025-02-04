# Receipt Processor API

## Setup

### Running Locally
```sh
git clone git@github.com:ariesjchang/receipt-processor-go.git
cd receipt-processor-go
go mod tidy
go run cmd/server/main.go
```

### Running with Docker
1. **Build the Docker image:**
```sh
docker build -t receipt-processor .
```

2. **Run the container:**
```sh
docker run -p 8080:8080 receipt-processor
```

## API Endpoints

### Process a Receipt
**Endpoint:** `POST /receipts/process`

**Request Body:**
```json
{
  "retailer": "Target",
  "purchaseDate": "2023-01-01",
  "purchaseTime": "14:01",
  "items": [
    {"shortDescription": "Item 1", "price": "5.00"},
    {"shortDescription": "Item 2", "price": "3.50"}
  ],
  "total": "8.50"
}
```

**Response:**
```json
{
  "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2"
}
```

### Get Receipt Points
**Endpoint:** `GET /receipts/{id}/points`

**Example Request:**
```sh
curl -X GET "http://localhost:8080/receipts/adb6b560-0eef-42bc-9d16-df48f30e89b2/points"
```

**Example Response:**
```json
{
  "points": 54
}
```

## Running Tests
To run the test suite, execute:
```sh
go test ./...
```

To run tests for a specific package:
```sh
go test ./internal/api/
```

To run tests with verbose output:
```sh
go test -v ./...
```
