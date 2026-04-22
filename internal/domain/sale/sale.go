package sale

import (
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"time"
)

type Sale struct {
	ID             string               `json:"id"`
	MemberID       string               `json:"member_id"`
	MemberName     string               `json:"member_name,omitempty"`
	UserID         string               `json:"user_id,omitempty"`
	Cashier        string               `json:"cashier,omitempty"`
	BranchID       string               `json:"branch_id,omitempty"`
	Payment        common.PaymentStatus `json:"payment"`
	Discount       int                  `json:"discount"`
	TotalSale      int                  `json:"total_sale"`
	ProfitEstimate int                  `json:"profit_estimate"`
	SaleDate       time.Time            `json:"sale_date"`
	CreatedAt      time.Time            `json:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
}

type Item struct {
	ID          string `json:"id"`
	SaleID      string `json:"sale_id,omitempty"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
	UnitName    string `json:"unit_name,omitempty"`
	Price       int    `json:"price"`
	Qty         int    `json:"qty"`
	SubTotal    int    `json:"sub_total"`
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
	Items []Sale
	Meta  ListMeta
}

type Detail struct {
	ID             string `json:"id"`
	MemberID       string `json:"member_id"`
	MemberName     string `json:"member_name"`
	SaleDate       string `json:"sale_date"`
	TotalSale      int    `json:"total_sale"`
	Discount       int    `json:"discount"`
	ProfitEstimate int    `json:"profit_estimate"`
	Payment        string `json:"payment"`
	Cashier        string `json:"cashier"`
	Items          []Item `json:"items"`
}

type UpdateRequest struct {
	MemberID *string `json:"member_id"`
	Discount *int    `json:"discount"`
	Payment  string  `json:"payment"`
}

type CreateSaleRequest struct {
	Sale struct {
		MemberID string `json:"member_id"`
		Payment  string `json:"payment"`
		Discount int    `json:"discount"`
	} `json:"sale" validate:"required"`
	SaleItems []struct {
		ProductID string `json:"product_id" validate:"required"`
		Price     int    `json:"price" validate:"required,min=1"`
		Qty       int    `json:"qty" validate:"required,min=1"`
	} `json:"sale_items" validate:"required,min=1,dive"`
}

type CreateItemRequest struct {
	SaleID    string `json:"sale_id" validate:"required"`
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}

type UpdateItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}
