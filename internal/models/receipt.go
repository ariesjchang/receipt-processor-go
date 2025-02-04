package models

type Receipt struct {
	ID           string
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []Item
	Total        string
	Points       int
}

type Item struct {
	ShortDescription string
	Price            string
}
