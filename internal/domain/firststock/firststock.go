package firststock

import (
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
)

type FirstStock struct {
	ID              string               `json:"id"`
	Description     string               `json:"description"`
	FirstStockDate  time.Time            `json:"first_stock_date"`
	BranchID        string               `json:"branch_id,omitempty"`
	UserID          string               `json:"user_id,omitempty"`
	TotalFirstStock int                  `json:"total_first_stock"`
	Payment         common.PaymentStatus `json:"payment"`
	CreatedAt       time.Time            `json:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"updated_at,omitempty"`
}

type Item struct {
	ID           string    `json:"id"`
	FirstStockID string    `json:"first_stock_id,omitempty"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name,omitempty"`
	UnitID       string    `json:"unit_id,omitempty"`
	UnitName     string    `json:"unit_name,omitempty"`
	Price        int       `json:"price"`
	Qty          int       `json:"qty"`
	SubTotal     int       `json:"sub_total"`
	ExpiredDate  time.Time `json:"expired_date,omitempty"`
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
	ID              string `json:"id"`
	Description     string `json:"description"`
	FirstStockDate  string `json:"first_stock_date"`
	TotalFirstStock int    `json:"total_first_stock"`
	Payment         string `json:"payment"`
}

type ListResult struct {
	Items []ListItem
	Meta  ListMeta
}

type Detail struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	FirstStockDate  string `json:"first_stock_date"`
	TotalFirstStock int    `json:"total_first_stock"`
	Payment         string `json:"payment"`
	Items           []Item `json:"items"`
}

type CreateRequest struct {
	Description    string `json:"description"`
	FirstStockDate string `json:"first_stock_date"`
}

type UpdateRequest struct {
	Description    string `json:"description"`
	FirstStockDate string `json:"first_stock_date"`
	Payment        string `json:"payment"`
}

type CreateItemRequest struct {
	FirstStockID string `json:"first_stock_id" validate:"required"`
	ProductID    string `json:"product_id" validate:"required"`
	UnitID       string `json:"unit_id" validate:"required"`
	Qty          int    `json:"qty" validate:"required,min=1"`
	ExpiredDate  string `json:"expired_date" validate:"required"`
}

type UpdateItemRequest struct {
	ProductID   string `json:"product_id" validate:"required"`
	UnitID      string `json:"unit_id" validate:"required"`
	Qty         int    `json:"qty" validate:"required,min=1"`
	ExpiredDate string `json:"expired_date" validate:"required"`
}
