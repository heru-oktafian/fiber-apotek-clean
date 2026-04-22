package product

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Products ports.ProductRepository
	IDs      ports.IDGenerator
}

func (s Service) Create(ctx context.Context, branchID string, input product.Product) (product.Product, error) {
	input.ID = s.IDs.New("PRD")
	input.BranchID = branchID
	input.Stock = 0
	if strings.TrimSpace(input.SKU) == "" {
		input.SKU = input.ID
	}
	if err := s.Products.Create(ctx, input); err != nil {
		return product.Product{}, apperror.New(http.StatusInternalServerError, "Failed to create resource", err)
	}
	return input, nil
}

func (s Service) List(ctx context.Context, branchID string, req product.ListRequest) (product.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Products.ListProducts(ctx, branchID, req)
	if err != nil {
		return product.ListResult{}, apperror.New(http.StatusInternalServerError, "Get all products failed", err)
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (product.Product, error) {
	item, err := s.Products.FindProductDetailByID(ctx, id, branchID)
	if err != nil {
		return product.Product{}, apperror.New(http.StatusNotFound, "Get product failed", err)
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, input product.Product) (product.Product, error) {
	existing, err := s.Products.FindProductDetailByID(ctx, id, branchID)
	if err != nil {
		return product.Product{}, apperror.New(http.StatusNotFound, "Update product failed", err)
	}
	input.ID = id
	input.BranchID = branchID
	input.Stock = existing.Stock
	if strings.TrimSpace(input.SKU) == "" {
		input.SKU = existing.SKU
	}
	if err := s.Products.Update(ctx, input); err != nil {
		return product.Product{}, apperror.New(http.StatusInternalServerError, "Update product failed", err)
	}
	return s.Products.FindProductDetailByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Products.FindProductDetailByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete product failed", err)
	}
	if err := s.Products.DeleteProduct(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete product failed", err)
	}
	return nil
}

func (s Service) SaleCombo(ctx context.Context, branchID, search string) ([]product.SaleComboItem, error) {
	items, err := s.Products.GetSaleCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get Combo Products failed", err)
	}
	return items, nil
}

func (s Service) PurchaseCombo(ctx context.Context, branchID, search string) ([]product.PurchaseComboItem, error) {
	items, err := s.Products.GetPurchaseCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get Combo Purchase Products failed", err)
	}
	return items, nil
}

func (s Service) OpnameCombo(ctx context.Context, branchID, search string) ([]product.OpnameComboItem, error) {
	items, err := s.Products.GetOpnameCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusNotFound, "Combobox tidak ditemukan", err)
	}
	return items, nil
}
