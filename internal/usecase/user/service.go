package user

import (
	"context"
	"net/http"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Users     ports.UserRepository
	Passwords ports.PasswordHasher
	IDs       ports.IDGenerator
}

func (s Service) List(ctx context.Context, req user.ListRequest) (user.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	result, err := s.Users.ListUsers(ctx, req)
	if err != nil {
		return user.ListResult{}, apperror.New(http.StatusInternalServerError, "Gagal mengambil data user", err.Error())
	}
	return result, nil
}

func (s Service) Detail(ctx context.Context, id string) (user.DetailWithBranches, error) {
	result, err := s.Users.FindUserWithBranches(ctx, id)
	if err != nil {
		return user.DetailWithBranches{}, apperror.New(http.StatusNotFound, "Pengguna tidak ditemukan", err.Error())
	}
	return result, nil
}

func (s Service) Create(ctx context.Context, req user.CreateRequest) (user.User, error) {
	if req.Username == "" || req.Name == "" || req.Password == "" || req.UserRole == "" {
		return user.User{}, apperror.New(http.StatusBadRequest, "Username, Password, Name dan Role harus diisi", nil)
	}
	if !isAllowedRole(req.UserRole) {
		return user.User{}, apperror.New(http.StatusBadRequest, "Invalid user role: "+req.UserRole, nil)
	}
	status := req.UserStatus
	if status == "" {
		status = "inactive"
	}
	if !isAllowedStatus(status) {
		return user.User{}, apperror.New(http.StatusBadRequest, "Invalid user status: "+status, nil)
	}
	hashed, err := s.Passwords.Hash(req.Password)
	if err != nil {
		return user.User{}, apperror.New(http.StatusInternalServerError, "Could not hash password", err.Error())
	}
	entity := user.User{
		ID:       s.IDs.New("USR"),
		Name:     req.Name,
		Username: req.Username,
		Password: hashed,
		Role:     common.UserRole(req.UserRole),
		Status:   status,
	}
	if err := s.Users.CreateUser(ctx, entity); err != nil {
		return user.User{}, apperror.New(http.StatusInternalServerError, "Gagal membuat user", err.Error())
	}
	entity.Password = ""
	return entity, nil
}

func (s Service) Update(ctx context.Context, id string, req user.UpdateRequest) (user.User, error) {
	entity, err := s.Users.FindByID(ctx, id)
	if err != nil {
		return user.User{}, apperror.New(http.StatusNotFound, "User tidak ditemukan", nil)
	}
	if req.Username != "" {
		entity.Username = req.Username
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.UserRole != "" {
		if !isAllowedRole(req.UserRole) {
			return user.User{}, apperror.New(http.StatusBadRequest, "Invalid user role: "+req.UserRole, nil)
		}
		entity.Role = common.UserRole(req.UserRole)
	}
	if req.UserStatus != "" {
		if !isAllowedStatus(req.UserStatus) {
			return user.User{}, apperror.New(http.StatusBadRequest, "Invalid user status: "+req.UserStatus, nil)
		}
		entity.Status = req.UserStatus
	}
	if req.Password != "" {
		hashed, err := s.Passwords.Hash(req.Password)
		if err != nil {
			return user.User{}, apperror.New(http.StatusInternalServerError, "Could not hash new password", err.Error())
		}
		entity.Password = hashed
	}
	if err := s.Users.UpdateUser(ctx, entity); err != nil {
		return user.User{}, apperror.New(http.StatusInternalServerError, "Gagal mengupdate user", err.Error())
	}
	entity.Password = ""
	return entity, nil
}

func isAllowedRole(role string) bool {
	allowed := map[string]bool{
		"administrator": true, "superadmin": true, "operator": true, "cashier": true,
		"finance": true, "pendaftaran": true, "rekammedis": true, "ralan": true,
		"ranap": true, "vk": true, "lab": true, "klaim": true, "simrs": true,
		"ipsrs": true, "umum": true,
	}
	return allowed[role]
}

func isAllowedStatus(status string) bool {
	allowed := map[string]bool{"active": true, "inactive": true}
	return allowed[status]
}
