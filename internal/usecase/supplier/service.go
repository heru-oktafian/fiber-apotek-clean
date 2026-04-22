package supplier

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/supplier"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Suppliers ports.SupplierRepository
	IDs       ports.IDGenerator
}

func (s Service) List(ctx context.Context, branchID string, req supplier.ListRequest) (supplier.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Suppliers.ListSuppliers(ctx, branchID, req)
	if err != nil {
		return supplier.ListResult{}, apperror.New(http.StatusInternalServerError, "Get suppliers failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (supplier.Supplier, error) {
	item, err := s.Suppliers.FindSupplierByID(ctx, id, branchID)
	if err != nil {
		return supplier.Supplier{}, apperror.New(http.StatusNotFound, "Get supplier failed", err.Error())
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, branchID string, req supplier.CreateRequest) (supplier.Supplier, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return supplier.Supplier{}, apperror.New(http.StatusBadRequest, "Create supplier failed", "name is required")
	}
	if req.SupplierCategoryID == 0 {
		return supplier.Supplier{}, apperror.New(http.StatusBadRequest, "Create supplier failed", "supplier_category_id is required")
	}
	item := supplier.Supplier{
		ID:                 s.IDs.New("SPL"),
		Name:               req.Name,
		Phone:              strings.TrimSpace(req.Phone),
		Address:            strings.TrimSpace(req.Address),
		PIC:                strings.TrimSpace(req.PIC),
		SupplierCategoryID: req.SupplierCategoryID,
		BranchID:           branchID,
	}
	if err := s.Suppliers.CreateSupplier(ctx, item); err != nil {
		return supplier.Supplier{}, apperror.New(http.StatusInternalServerError, "Create supplier failed", err.Error())
	}
	return s.Suppliers.FindSupplierByID(ctx, item.ID, branchID)
}

func (s Service) Combo(ctx context.Context, branchID, search string) ([]supplier.ComboItem, error) {
	items, err := s.Suppliers.GetSupplierCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get supplier combo failed", err.Error())
	}
	return items, nil
}
