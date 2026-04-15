package opname

import "time"

type Opname struct {
	ID          string
	Description string
	BranchID    string
	UserID      string
	OpnameDate  time.Time
	TotalOpname int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Item struct {
	ID            string
	OpnameID      string
	ProductID     string
	ProductName   string
	Qty           int
	QtyExist      int
	Price         int
	SubTotal      int
	SubTotalExist int
	ExpiredDate   time.Time
}

type Detail struct {
	ID          string
	Description string
	OpnameDate  time.Time
	TotalOpname int
	Items       []Item
}

type CreateOpnameRequest struct {
	Description string `json:"description"`
	OpnameDate  string `json:"opname_date" validate:"required"`
}

type CreateOpnameItemRequest struct {
	OpnameID    string `json:"opname_id" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	Qty         int    `json:"qty" validate:"required,min=0"`
	Price       int    `json:"price" validate:"required,min=1"`
	ExpiredDate string `json:"expired_date" validate:"required"`
}
