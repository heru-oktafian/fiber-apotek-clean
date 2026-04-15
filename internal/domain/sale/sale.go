package sale

import (
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"time"
)

type Sale struct {
	ID             string
	MemberID       string
	UserID         string
	BranchID       string
	Payment        common.PaymentStatus
	Discount       int
	TotalSale      int
	ProfitEstimate int
	SaleDate       time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Item struct {
	ID        string
	SaleID    string
	ProductID string
	Price     int
	Qty       int
	SubTotal  int
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
