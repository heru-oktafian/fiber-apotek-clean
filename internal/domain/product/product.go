package product

import "time"

type Product struct {
	ID                string
	SKU               string
	Name              string
	Description       string
	BranchID          string
	UnitID            string
	Stock             int
	PurchasePrice     int
	SalesPrice        int
	AlternatePrice    int
	ProductCategoryID string
	ExpiredDate       time.Time
}

type SaleComboItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	UnitID      string `json:"unit_id"`
	UnitName    string `json:"unit_name"`
}

type PurchaseComboItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
	UnitID      string `json:"unit_id"`
	UnitName    string `json:"unit_name"`
}

type OpnameComboItem struct {
	ProductID   string `json:"pro_id"`
	ProductName string `json:"pro_name"`
	UnitID      string `json:"unit_id"`
	Stock       int    `json:"stock"`
	UnitName    string `json:"unit_name"`
	Price       int    `json:"price"`
}
