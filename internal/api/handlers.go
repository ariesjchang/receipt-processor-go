package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
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
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// Validate required fields
	var totalPattern = regexp.MustCompile(`^\d+\.\d{2}$`)
	if receipt.Retailer == "" || receipt.Total == "" || receipt.PurchaseDate == "" ||
		receipt.PurchaseTime == "" || len(receipt.Items) == 0 || !totalPattern.MatchString(receipt.Total) {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
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
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	mu.RLock()
	receipt, exists := receipts[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}
