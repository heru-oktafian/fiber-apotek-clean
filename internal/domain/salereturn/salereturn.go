package salereturn

import "time"

type SaleReturn struct {
	ID          string    `json:"id"`
	SaleID      string    `json:"sale_id"`
	ReturnDate  time.Time `json:"return_date"`
	BranchID    string    `json:"branch_id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Payment     string    `json:"payment"`
	TotalReturn int       `json:"total_return"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Item struct {
	ID           string    `json:"id"`
	SaleReturnID string    `json:"sale_return_id,omitempty"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name,omitempty"`
	UnitID       string    `json:"unit_id,omitempty"`
	UnitName     string    `json:"unit_name,omitempty"`
	Price        int       `json:"price"`
	Qty          int       `json:"qty"`
	SubTotal     int       `json:"sub_total"`
	ExpiredDate  time.Time `json:"expired_date"`
}

type ListRequest struct {
	Search string
	Month  string
	Page   int
	Limit  int
}

type ListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	Month     string `json:"month"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type ListItem struct {
	ID          string `json:"id"`
	SaleID      string `json:"sale_id"`
	ReturnDate  string `json:"return_date"`
	TotalReturn int    `json:"total_purchase"`
	Payment     string `json:"payment"`
}

type ListResult struct {
	Items []ListItem
	Meta  ListMeta
}

type Detail struct {
	ID          string          `json:"id"`
	SaleID      string          `json:"sale_id"`
	ReturnDate  string          `json:"return_date"`
	TotalReturn int             `json:"total_purchase"`
	Payment     string          `json:"payment"`
	Items       []FormattedItem `json:"items"`
}

type FormattedItem struct {
	ID           string `json:"id"`
	SaleReturnID string `json:"sale_return_id"`
	ProductID    string `json:"pro_id"`
	ProductName  string `json:"pro_name"`
	UnitID       string `json:"unit_id"`
	UnitName     string `json:"unit_name"`
	Price        int    `json:"price"`
	Qty          int    `json:"qty"`
	SubTotal     int    `json:"sub_total"`
	ExpiredDate  string `json:"expired_date"`
}

type CreateRequest struct {
	SaleReturn struct {
		SaleID     string `json:"sale_id" validate:"required"`
		ReturnDate string `json:"return_date"`
		Payment    string `json:"payment"`
	} `json:"sale_return" validate:"required"`
	SaleReturnItems []struct {
		ProductID   string `json:"product_id" validate:"required"`
		Qty         int    `json:"qty" validate:"required,min=1"`
		ExpiredDate string `json:"expired_date" validate:"required"`
	} `json:"sale_return_items" validate:"required,min=1,dive"`
}

type SaleComboItem struct {
	ID        string    `json:"id"`
	SaleDate  time.Time `json:"sale_date"`
	TotalSale int       `json:"total_sale"`
	MemberID  string    `json:"member_id"`
}

type ReturnableItem struct {
	ProductID   string `json:"pro_id"`
	ProductName string `json:"pro_name"`
	Stock       int    `json:"stock"`
	UnitID      string `json:"unit_id"`
	UnitName    string `json:"unit_name"`
	Price       int    `json:"price"`
}
