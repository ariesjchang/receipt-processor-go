package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/ariesjchang/receipt-processor-go/internal/models"
	"github.com/ariesjchang/receipt-processor-go/internal/services"
)

var (
	receipts = make(map[string]models.Receipt)
	mu       sync.RWMutex
)

// ProcessReceipt handles the receipt processing endpoint
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if receipt.Retailer == "" || receipt.Total == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || receipt.Items == nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate ID and calculate points
	receipt.ID = uuid.New().String()
	receipt.Points = services.CalculatePoints(receipt)

	// Store receipt
	mu.Lock()
	receipts[receipt.ID] = receipt
	mu.Unlock()

	// Respond with the receipt ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receipt.ID})
}

// GetPoints handles fetching the points for a receipt by ID
func GetPoints(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/receipts/"):]
	id = id[:len(id)-len("/points")]

	mu.RLock()
	receipt, exists := receipts[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}
