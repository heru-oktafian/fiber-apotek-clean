package duplicatereceipt

import (
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"time"
)

type DuplicateReceipt struct {
	ID                    string               `json:"id"`
	MemberID              string               `json:"member_id"`
	MemberName            string               `json:"member_name,omitempty"`
	Description           string               `json:"description"`
	DuplicateReceiptDate  time.Time            `json:"duplicate_receipt_date"`
	TotalDuplicateReceipt int                  `json:"total_duplicate_receipt"`
	ProfitEstimate        int                  `json:"profit_estimate"`
	Payment               common.PaymentStatus `json:"payment"`
	BranchID              string               `json:"branch_id,omitempty"`
	UserID                string               `json:"user_id,omitempty"`
	CreatedAt             time.Time            `json:"created_at,omitempty"`
	UpdatedAt             time.Time            `json:"updated_at,omitempty"`
}

type Item struct {
	ID                 string `json:"id"`
	DuplicateReceiptID string `json:"duplicate_receipt_id"`
	ProductID          string `json:"product_id"`
	ProductName        string `json:"product_name,omitempty"`
	UnitName           string `json:"unit_name,omitempty"`
	Price              int    `json:"price"`
	Qty                int    `json:"qty"`
	SubTotal           int    `json:"sub_total"`
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
	ID                    string `json:"id"`
	MemberID              string `json:"member_id"`
	MemberName            string `json:"member_name"`
	DuplicateReceiptDate  string `json:"duplicate_receipt_date"`
	TotalDuplicateReceipt int    `json:"total_duplicate_receipt"`
	ProfitEstimate        int    `json:"profit_estimate"`
	Payment               string `json:"payment"`
}

type ListResult struct {
	Items []ListItem
	Meta  ListMeta
}

type DetailSummaryItem struct {
	ID                    string `json:"id"`
	DuplicateReceiptDate  string `json:"duplicate_receipt_date"`
	Description           string `json:"description"`
	Payment               string `json:"payment"`
	TotalDuplicateReceipt int    `json:"total_duplicate_receipt"`
}

type DetailSummaryResult struct {
	Items []DetailSummaryItem
	Meta  ListMeta
}

type Detail struct {
	ID                    string `json:"id"`
	MemberID              string `json:"member_id"`
	MemberName            string `json:"member_name"`
	Description           string `json:"description"`
	DuplicateReceiptDate  string `json:"duplicate_receipt_date"`
	TotalDuplicateReceipt int    `json:"total_duplicate_receipt"`
	ProfitEstimate        int    `json:"profit_estimate"`
	Payment               string `json:"payment"`
}

type CreateRequest struct {
	DuplicateReceipt struct {
		MemberID             string `json:"member_id"`
		Description          string `json:"description"`
		DuplicateReceiptDate string `json:"duplicate_receipt_date" validate:"required"`
		Payment              string `json:"payment"`
	} `json:"duplicate_receipt" validate:"required"`
	Items []struct {
		ProductID string `json:"product_id" validate:"required"`
		Price     int    `json:"price"`
		Qty       int    `json:"qty" validate:"required,min=1"`
		SubTotal  int    `json:"sub_total"`
	} `json:"items" validate:"required,min=1,dive"`
}

type UpdateRequest struct {
	MemberID    *string `json:"member_id"`
	Description *string `json:"description"`
	Payment     string  `json:"payment"`
}

type CreateItemRequest struct {
	DuplicateReceiptID string `json:"duplicate_receipt_id" validate:"required"`
	ProductID          string `json:"product_id" validate:"required"`
	Qty                int    `json:"qty" validate:"required,min=1"`
}

type UpdateItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=1"`
}
