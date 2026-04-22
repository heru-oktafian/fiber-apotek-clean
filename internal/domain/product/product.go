package product

import "time"

type Product struct {
	ID                  string    `json:"id"`
	SKU                 string    `json:"sku"`
	Name                string    `json:"name"`
	Alias               string    `json:"alias"`
	Description         string    `json:"description"`
	Ingredient          string    `json:"ingredient"`
	Dosage              string    `json:"dosage"`
	SideAffection       string    `json:"side_affection"`
	BranchID            string    `json:"branch_id,omitempty"`
	UnitID              string    `json:"unit_id"`
	UnitName            string    `json:"unit_name,omitempty"`
	Stock               int       `json:"stock"`
	PurchasePrice       int       `json:"purchase_price"`
	SalesPrice          int       `json:"sales_price"`
	AlternatePrice      int       `json:"alternate_price"`
	ProductCategoryID   uint      `json:"product_category_id"`
	ProductCategoryName string    `json:"product_category_name,omitempty"`
	ExpiredDate         time.Time `json:"expired_date"`
}

type ListRequest struct {
	Search string
	Page   int
	Limit  int
}

type ListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type ListResult struct {
	Items []Product
	Meta  ListMeta
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
