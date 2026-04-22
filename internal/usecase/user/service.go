package user

import (
	"context"
	"net/http"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Users ports.UserRepository
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
