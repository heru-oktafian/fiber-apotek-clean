package supplier

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/supplier"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) Update(ctx context.Context, branchID, id string, req supplier.CreateRequest) (supplier.Supplier, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return supplier.Supplier{}, apperror.New(http.StatusBadRequest, "Update supplier failed", "name is required")
	}
	if req.SupplierCategoryID == 0 {
		return supplier.Supplier{}, apperror.New(http.StatusBadRequest, "Update supplier failed", "supplier_category_id is required")
	}
	if _, err := s.Suppliers.FindSupplierByID(ctx, id, branchID); err != nil {
		return supplier.Supplier{}, apperror.New(http.StatusNotFound, "Update supplier failed", "supplier not found")
	}
	item := supplier.Supplier{ID: id, Name: req.Name, Phone: strings.TrimSpace(req.Phone), Address: strings.TrimSpace(req.Address), PIC: strings.TrimSpace(req.PIC), SupplierCategoryID: req.SupplierCategoryID, BranchID: branchID}
	if err := s.Suppliers.UpdateSupplier(ctx, item); err != nil {
		return supplier.Supplier{}, apperror.New(http.StatusInternalServerError, "Update supplier failed", err.Error())
	}
	return s.Suppliers.FindSupplierByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Suppliers.FindSupplierByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete supplier failed", "supplier not found")
	}
	if err := s.Suppliers.DeleteSupplier(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete supplier failed", err.Error())
	}
	return nil
}

func (s Service) Combo(ctx context.Context, branchID, search string) ([]supplier.ComboItem, error) {
	items, err := s.Suppliers.GetSupplierCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get supplier combo failed", err.Error())
	}
	return items, nil
}

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Suppliers.ListSuppliers(ctx, branchID, supplier.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export suppliers excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Suppliers")
	sheet := "Suppliers"
	f.SetCellValue(sheet, "A1", "DATA SUPPLIERS")
	headers := []string{"ID", "NAME", "PHONE", "ADDRESS", "PIC", "CATEGORY"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Address)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.PIC)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.SupplierCategory)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export suppliers excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("suppliers-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Suppliers.ListSuppliers(ctx, branchID, supplier.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export suppliers pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER SUPPLIERS")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER SUPPLIERS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "NAME", "PHONE", "ADDRESS", "PIC", "CATEGORY"}
	widths := []float64{35, 50, 35, 70, 35, 52}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items.Items {
		values := []string{item.ID, item.Name, item.Phone, item.Address, item.PIC, item.SupplierCategory}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export suppliers pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("suppliers-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
