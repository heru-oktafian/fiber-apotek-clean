package membercategory

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/membercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Categories ports.MemberCategoryRepository
}

func (s Service) List(ctx context.Context, branchID string, req membercategory.ListRequest) (membercategory.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Categories.ListMemberCategories(ctx, branchID, req)
	if err != nil {
		return membercategory.ListResult{}, apperror.New(http.StatusInternalServerError, "Get member categories failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID string, id uint) (membercategory.MemberCategory, error) {
	item, err := s.Categories.FindMemberCategoryByID(ctx, id, branchID)
	if err != nil {
		return membercategory.MemberCategory{}, apperror.New(http.StatusNotFound, "Get member category failed", err.Error())
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, branchID string, req membercategory.CreateRequest) (membercategory.MemberCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return membercategory.MemberCategory{}, apperror.New(http.StatusBadRequest, "Create member category failed", "name is required")
	}
	item, err := s.Categories.CreateMemberCategory(ctx, membercategory.MemberCategory{Name: name, PointsConversionRate: req.PointsConversionRate, BranchID: branchID})
	if err != nil {
		return membercategory.MemberCategory{}, apperror.New(http.StatusInternalServerError, "Create member category failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID string, id uint, req membercategory.CreateRequest) (membercategory.MemberCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return membercategory.MemberCategory{}, apperror.New(http.StatusBadRequest, "Update member category failed", "name is required")
	}
	if _, err := s.Categories.FindMemberCategoryByID(ctx, id, branchID); err != nil {
		return membercategory.MemberCategory{}, apperror.New(http.StatusNotFound, "Update member category failed", "member category not found")
	}
	if err := s.Categories.UpdateMemberCategory(ctx, membercategory.MemberCategory{ID: id, Name: name, PointsConversionRate: req.PointsConversionRate, BranchID: branchID}); err != nil {
		return membercategory.MemberCategory{}, apperror.New(http.StatusInternalServerError, "Update member category failed", err.Error())
	}
	return s.Categories.FindMemberCategoryByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID string, id uint) error {
	if _, err := s.Categories.FindMemberCategoryByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete member category failed", "member category not found")
	}
	if err := s.Categories.DeleteMemberCategory(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete member category failed", err.Error())
	}
	return nil
}

func (s Service) Combo(ctx context.Context, branchID, search string) ([]membercategory.ComboItem, error) {
	items, err := s.Categories.GetMemberCategoryCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get member category combo failed", err.Error())
	}
	return items, nil
}
