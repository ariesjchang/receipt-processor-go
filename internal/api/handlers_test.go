package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ariesjchang/receipt-processor-go/internal/api"
	"github.com/ariesjchang/receipt-processor-go/internal/models"
)

func TestProcessReceipt_ValidRequest(t *testing.T) {
	reqBody := `{
		"retailer": "Target",
		"purchaseDate": "2023-01-01",
		"purchaseTime": "14:01",
		"items": [
			{"shortDescription": "Item 1", "price": "5.00"},
			{"shortDescription": "Item 2", "price": "3.50"}
		],
		"total": "8.50"
	}`

	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.ProcessReceipt(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
}

func TestProcessReceipt_InvalidJSON(t *testing.T) {
	reqBody := `{"retailer": "Target", "purchaseDate": "2023-01-01", "purchaseTime": "14:01", "items": [`
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.ProcessReceipt(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

func TestProcessReceipt_MissingFields(t *testing.T) {
	reqBody := `{"purchaseDate": "2023-01-01"}`
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.ProcessReceipt(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

func TestGetPoints_ValidReceipt(t *testing.T) {
	// Step 1: Create a receipt
	receipt := models.Receipt{
		Retailer:     "Test Store",
		Total:        "10.00",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "14:01",
		Items:        []models.Item{{ShortDescription: "Item 1", Price: "5.00"}},
	}
	body, _ := json.Marshal(receipt)

	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.ProcessReceipt(w, req)

	// Step 2: Extract ID from response
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	receiptID, exists := resp["id"]
	if !exists {
		t.Fatalf("Expected receipt ID in response, got: %v", resp)
	}

	// Step 3: Fetch points using the generated ID
	req = httptest.NewRequest(http.MethodGet, "/receipts/"+receiptID+"/points", nil)
	w = httptest.NewRecorder()
	api.GetPoints(w, req)

	// Step 4: Verify response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
}

func TestGetPoints_InvalidReceipt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/receipts/invalid/points", nil)
	w := httptest.NewRecorder()

	api.GetPoints(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}
