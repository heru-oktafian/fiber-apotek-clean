package purchase

import (
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"time"
)

type Purchase struct {
	ID            string               `json:"id"`
	SupplierID    string               `json:"supplier_id"`
	SupplierName  string               `json:"supplier_name,omitempty"`
	PurchaseDate  time.Time            `json:"purchase_date"`
	BranchID      string               `json:"branch_id,omitempty"`
	UserID        string               `json:"user_id,omitempty"`
	Payment       common.PaymentStatus `json:"payment"`
	TotalPurchase int                  `json:"total_purchase"`
	CreatedAt     time.Time            `json:"created_at,omitempty"`
	UpdatedAt     time.Time            `json:"updated_at,omitempty"`
}

type Item struct {
	ID          string    `json:"id"`
	PurchaseID  string    `json:"purchase_id,omitempty"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	UnitID      string    `json:"unit_id"`
	UnitName    string    `json:"unit_name,omitempty"`
	Price       int       `json:"price"`
	Qty         int       `json:"qty"`
	SubTotal    int       `json:"sub_total"`
	ExpiredDate time.Time `json:"expired_date"`
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
	Items []Purchase
	Meta  ListMeta
}

type Detail struct {
	ID            string          `json:"id"`
	SupplierID    string          `json:"supplier_id"`
	SupplierName  string          `json:"supplier_name"`
	PurchaseDate  string          `json:"purchase_date"`
	TotalPurchase int             `json:"total_purchase"`
	Payment       string          `json:"payment"`
	Items         []FormattedItem `json:"items"`
}

type FormattedItem struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	UnitID      string `json:"unit_id"`
	UnitName    string `json:"unit_name"`
	Price       int    `json:"price"`
	Qty         int    `json:"qty"`
	SubTotal    int    `json:"sub_total"`
	ExpiredDate string `json:"expired_date"`
}

type UpdateRequest struct {
	SupplierID   string `json:"supplier_id"`
	PurchaseDate string `json:"purchase_date"`
	Payment      string `json:"payment"`
}

type CreatePurchaseRequest struct {
	Purchase struct {
		SupplierID   string `json:"supplier_id" validate:"required"`
		PurchaseDate string `json:"purchase_date"`
		Payment      string `json:"payment"`
	} `json:"purchase" validate:"required"`
	PurchaseItems []struct {
		ProductID   string `json:"product_id" validate:"required"`
		UnitID      string `json:"unit_id" validate:"required"`
		Price       int    `json:"price" validate:"required,min=1"`
		Qty         int    `json:"qty" validate:"required,min=1"`
		ExpiredDate string `json:"expired_date" validate:"required"`
	} `json:"purchase_items" validate:"required,min=1,dive"`
}

type CreateItemRequest struct {
	PurchaseID  string `json:"purchase_id" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	UnitID      string `json:"unit_id" validate:"required"`
	Price       int    `json:"price" validate:"required,min=1"`
	Qty         int    `json:"qty" validate:"required,min=1"`
	ExpiredDate string `json:"expired_date" validate:"required"`
}

type UpdateItemRequest struct {
	ProductID   string `json:"product_id" validate:"required"`
	UnitID      string `json:"unit_id" validate:"required"`
	Price       int    `json:"price" validate:"required,min=1"`
	Qty         int    `json:"qty" validate:"required,min=1"`
	ExpiredDate string `json:"expired_date" validate:"required"`
}
