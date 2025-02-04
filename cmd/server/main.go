package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ariesjchang/receipt-processor-go/internal/api"
)

func main() {
	fmt.Println("Starting server on :8080...")
	http.HandleFunc("/receipts/process", api.ProcessReceipt)
	http.HandleFunc("/receipts/", api.GetPoints)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
