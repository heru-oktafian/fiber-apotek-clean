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
