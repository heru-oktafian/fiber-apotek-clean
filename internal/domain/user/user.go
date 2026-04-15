package user

import "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"

type User struct {
	ID       string
	Name     string
	Username string
	Password string
	Role     common.UserRole
	Status   string
}
