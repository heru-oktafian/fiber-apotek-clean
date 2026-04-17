package branch

import (
	"context"
	"net/http"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Branches ports.BranchRepository
}

func (s Service) List(ctx context.Context, req branch.ListRequest) (branch.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Branches.ListBranches(ctx, req)
	if err != nil {
		return branch.ListResult{}, apperror.New(http.StatusInternalServerError, "Get all branch failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, id string) (branch.Branch, error) {
	item, err := s.Branches.FindBranchByID(ctx, id)
	if err != nil {
		return branch.Branch{}, apperror.New(http.StatusNotFound, "Get branch failed", err.Error())
	}
	return item, nil
}
