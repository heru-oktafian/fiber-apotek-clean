package suppliercategory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/suppliercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListSupplierCategories(ctx, branchID, suppliercategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export supplier categories excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Supplier Categories")
	sheet := "Supplier Categories"
	f.SetCellValue(sheet, "A1", "DATA SUPPLIER CATEGORIES")
	f.SetCellValue(sheet, "A3", "ID")
	f.SetCellValue(sheet, "B3", "NAME")
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Name)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export supplier categories excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("supplier-categories-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListSupplierCategories(ctx, branchID, suppliercategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export supplier categories pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER SUPPLIER CATEGORIES")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER SUPPLIER CATEGORIES", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 8, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(217, 8, "NAME", "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	for _, item := range items.Items {
		pdf.CellFormat(60, 8, fmt.Sprintf("%d", item.ID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(217, 8, item.Name, "1", 1, "L", false, 0, "")
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export supplier categories pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("supplier-categories-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
