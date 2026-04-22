package userbranch

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/userbranch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Branches ports.BranchRepository
	Users    ports.UserRepository
}

func (s Service) List(ctx context.Context) ([]userbranch.Detail, error) {
	items, err := s.Branches.ListAllUserBranches(ctx)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get user branches failed", "Failed to fetch user branches with details")
	}
	return items, nil
}

func (s Service) GetByKeys(ctx context.Context, userID, branchID string) ([]userbranch.Detail, error) {
	items, err := s.Branches.FindUserBranchDetail(ctx, userID, branchID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}
	return items, nil
}

func (s Service) Create(ctx context.Context, req userbranch.CreateRequest) error {
	req.UserID = strings.TrimSpace(req.UserID)
	req.BranchID = strings.TrimSpace(req.BranchID)
	if req.UserID == "" || req.BranchID == "" {
		return apperror.New(http.StatusBadRequest, "Create userbranch failed", "user_id and branch_id are required")
	}
	if _, err := s.Users.FindByID(ctx, req.UserID); err != nil {
		return apperror.New(http.StatusNotFound, "Create userbranch failed", "user not found")
	}
	if _, err := s.Branches.FindBranchByID(ctx, req.BranchID); err != nil {
		return apperror.New(http.StatusNotFound, "Create userbranch failed", "branch not found")
	}
	exists, err := s.Branches.UserHasBranch(ctx, req.UserID, req.BranchID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Create userbranch failed", err.Error())
	}
	if exists {
		return apperror.New(http.StatusBadRequest, "Create userbranch failed", "userbranch already exists")
	}
	if err := s.Users.CreateUserBranch(ctx, req.UserID, req.BranchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Create userbranch failed", err.Error())
	}
	return nil
}
