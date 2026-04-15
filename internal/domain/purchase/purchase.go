package purchase

import (
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"time"
)

type Purchase struct {
	ID            string
	SupplierID    string
	PurchaseDate  time.Time
	BranchID      string
	UserID        string
	Payment       common.PaymentStatus
	TotalPurchase int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Item struct {
	ID          string
	PurchaseID  string
	ProductID   string
	UnitID      string
	Price       int
	Qty         int
	SubTotal    int
	ExpiredDate time.Time
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
