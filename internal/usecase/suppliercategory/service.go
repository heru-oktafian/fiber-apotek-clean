package suppliercategory

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/suppliercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Categories ports.SupplierCategoryRepository
}

func (s Service) List(ctx context.Context, branchID string, req suppliercategory.ListRequest) (suppliercategory.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Categories.ListSupplierCategories(ctx, branchID, req)
	if err != nil {
		return suppliercategory.ListResult{}, apperror.New(http.StatusInternalServerError, "Get supplier categories failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID string, id uint) (suppliercategory.SupplierCategory, error) {
	item, err := s.Categories.FindSupplierCategoryByID(ctx, id, branchID)
	if err != nil {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusNotFound, "Get supplier category failed", err.Error())
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, branchID string, req suppliercategory.CreateRequest) (suppliercategory.SupplierCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusBadRequest, "Create supplier category failed", "name is required")
	}
	item, err := s.Categories.CreateSupplierCategory(ctx, suppliercategory.SupplierCategory{Name: name, BranchID: branchID})
	if err != nil {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusInternalServerError, "Create supplier category failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID string, id uint, req suppliercategory.CreateRequest) (suppliercategory.SupplierCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusBadRequest, "Update supplier category failed", "name is required")
	}
	if _, err := s.Categories.FindSupplierCategoryByID(ctx, id, branchID); err != nil {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusNotFound, "Update supplier category failed", "supplier category not found")
	}
	if err := s.Categories.UpdateSupplierCategory(ctx, suppliercategory.SupplierCategory{ID: id, Name: name, BranchID: branchID}); err != nil {
		return suppliercategory.SupplierCategory{}, apperror.New(http.StatusInternalServerError, "Update supplier category failed", err.Error())
	}
	return s.Categories.FindSupplierCategoryByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID string, id uint) error {
	if _, err := s.Categories.FindSupplierCategoryByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete supplier category failed", "supplier category not found")
	}
	if err := s.Categories.DeleteSupplierCategory(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete supplier category failed", err.Error())
	}
	return nil
}

func (s Service) Combo(ctx context.Context, branchID string) ([]suppliercategory.ComboItem, error) {
	items, err := s.Categories.GetSupplierCategoryCombo(ctx, branchID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get supplier category combo failed", err.Error())
	}
	return items, nil
}
