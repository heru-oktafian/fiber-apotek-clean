package user

import "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"

type User struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Username string          `json:"username"`
	Password string          `json:"password,omitempty"`
	Role     common.UserRole `json:"user_role"`
	Status   string          `json:"user_status"`
}

type ListRequest struct {
	Search string
	Page   int
	Limit  int
}

type ListResult struct {
	Items []User
	Meta  ListMeta
}

type ListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type DetailWithBranches struct {
	User           User           `json:"user"`
	DetailBranches []BranchDetail `json:"detail_branches"`
}

type BranchDetail struct {
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}
