package anotherincome

import (
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
)

type AnotherIncome struct {
	ID          string               `json:"id"`
	Description string               `json:"description"`
	IncomeDate  time.Time            `json:"income_date"`
	BranchID    string               `json:"branch_id,omitempty"`
	UserID      string               `json:"user_id,omitempty"`
	TotalIncome int                  `json:"total_income"`
	Payment     common.PaymentStatus `json:"payment"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
}

type CreateRequest struct {
	IncomeDate  string `json:"income_date" validate:"required"`
	Description string `json:"description"`
	TotalIncome int    `json:"total_income" validate:"required,min=0"`
	Payment     string `json:"payment"`
}

type UpdateRequest struct {
	IncomeDate  string `json:"income_date" validate:"required"`
	Description string `json:"description"`
	TotalIncome int    `json:"total_income" validate:"required,min=0"`
	Payment     string `json:"payment"`
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
	Description string `json:"description"`
	IncomeDate  string `json:"income_date"`
	TotalIncome int    `json:"total_income"`
	Payment     string `json:"payment"`
}

type ListResult struct {
	Items []ListItem
	Meta  ListMeta
}
