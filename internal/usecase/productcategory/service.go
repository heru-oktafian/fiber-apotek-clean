package productcategory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/productcategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type Service struct {
	Categories ports.ProductCategoryRepository
}

func (s Service) List(ctx context.Context, branchID string, req productcategory.ListRequest) (productcategory.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Categories.ListProductCategories(ctx, branchID, req)
	if err != nil {
		return productcategory.ListResult{}, apperror.New(http.StatusInternalServerError, "Get product categories failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID string, id uint) (productcategory.ProductCategory, error) {
	item, err := s.Categories.FindProductCategoryByID(ctx, id, branchID)
	if err != nil {
		return productcategory.ProductCategory{}, apperror.New(http.StatusNotFound, "Get product category failed", err.Error())
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, branchID string, req productcategory.CreateRequest) (productcategory.ProductCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return productcategory.ProductCategory{}, apperror.New(http.StatusBadRequest, "Create product category failed", "name is required")
	}
	item, err := s.Categories.CreateProductCategory(ctx, productcategory.ProductCategory{Name: name, BranchID: branchID})
	if err != nil {
		return productcategory.ProductCategory{}, apperror.New(http.StatusInternalServerError, "Create product category failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID string, id uint, req productcategory.CreateRequest) (productcategory.ProductCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return productcategory.ProductCategory{}, apperror.New(http.StatusBadRequest, "Update product category failed", "name is required")
	}
	if _, err := s.Categories.FindProductCategoryByID(ctx, id, branchID); err != nil {
		return productcategory.ProductCategory{}, apperror.New(http.StatusNotFound, "Update product category failed", "product category not found")
	}
	if err := s.Categories.UpdateProductCategory(ctx, productcategory.ProductCategory{ID: id, Name: name, BranchID: branchID}); err != nil {
		return productcategory.ProductCategory{}, apperror.New(http.StatusInternalServerError, "Update product category failed", err.Error())
	}
	return s.Categories.FindProductCategoryByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID string, id uint) error {
	if _, err := s.Categories.FindProductCategoryByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete product category failed", "product category not found")
	}
	if err := s.Categories.DeleteProductCategory(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete product category failed", err.Error())
	}
	return nil
}

func (s Service) Combo(ctx context.Context, branchID, search string) ([]productcategory.ComboItem, error) {
	items, err := s.Categories.GetProductCategoryCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get product category combo failed", err.Error())
	}
	return items, nil
}

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListProductCategories(ctx, branchID, productcategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export product categories excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Product Categories")
	sheet := "Product Categories"
	f.SetCellValue(sheet, "A1", "DATA PRODUCT CATEGORIES")
	f.SetCellValue(sheet, "A3", "ID")
	f.SetCellValue(sheet, "B3", "NAME")
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ProductCategoryID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.ProductCategoryName)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export product categories excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("product-categories-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListProductCategories(ctx, branchID, productcategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export product categories pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER PRODUCT CATEGORIES")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER PRODUCT CATEGORIES", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 8, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(217, 8, "NAME", "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	for _, item := range items.Items {
		pdf.CellFormat(60, 8, fmt.Sprintf("%d", item.ProductCategoryID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(217, 8, item.ProductCategoryName, "1", 1, "L", false, 0, "")
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export product categories pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("product-categories-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
