package userbranch

import (
	"context"
	"net/http"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/userbranch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Branches ports.BranchRepository
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
