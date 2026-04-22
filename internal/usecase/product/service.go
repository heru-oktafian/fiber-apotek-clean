package product

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Products.ListProducts(ctx, branchID, product.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export products excel failed", err)
	}
	f := exportshared.NewExcelFile("Produk")
	sheet := "Produk"
	f.SetCellValue(sheet, "A1", "DATA PRODUK")
	headers := []string{"ID", "SKU", "NAME", "ALIAS", "PURCHASE PRI", "SALE PRI", "ALTERNATIF PRI", "STOCK", "UNIT", "EXPIRED DATE"}
	for i, h := range headers {
		cell, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", cell), h)
	}
	for i, p := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.SKU)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Alias)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), p.PurchasePrice)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.SalesPrice)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.AlternatePrice)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("%d %s", p.Stock, p.UnitName))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), p.UnitName)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), p.ExpiredDate.Format("02/01/2006"))
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export products excel failed", err)
	}
	return bytes, fmt.Sprintf("products-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Products.ListProducts(ctx, branchID, product.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export products pdf failed", err)
	}
	pdf := exportshared.NewPDF("MASTER PRODUCTS")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER PRODUCTS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"SKU", "NAME", "ALIAS", "PURCHASE", "SALE", "ALT", "STOCK", "UNIT", "EXPIRED"}
	widths := []float64{28, 40, 35, 25, 25, 25, 20, 20, 30}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, p := range items.Items {
		values := []string{p.SKU, p.Name, p.Alias, fmt.Sprintf("%d", p.PurchasePrice), fmt.Sprintf("%d", p.SalesPrice), fmt.Sprintf("%d", p.AlternatePrice), fmt.Sprintf("%d", p.Stock), p.UnitName, p.ExpiredDate.Format("02/01/2006")}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export products pdf failed", err)
	}
	return bytes, fmt.Sprintf("products-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
