package membercategory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/membercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListMemberCategories(ctx, branchID, membercategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export member categories excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Member Categories")
	sheet := "Member Categories"
	f.SetCellValue(sheet, "A1", "DATA MEMBER CATEGORIES")
	f.SetCellValue(sheet, "A3", "ID")
	f.SetCellValue(sheet, "B3", "NAME")
	f.SetCellValue(sheet, "C3", "POINTS CONVERSION RATE")
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.PointsConversionRate)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export member categories excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("member-categories-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Categories.ListMemberCategories(ctx, branchID, membercategory.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export member categories pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER MEMBER CATEGORIES")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER MEMBER CATEGORIES", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 8, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(150, 8, "NAME", "1", 0, "C", false, 0, "")
	pdf.CellFormat(87, 8, "POINTS CONVERSION RATE", "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	for _, item := range items.Items {
		pdf.CellFormat(40, 8, fmt.Sprintf("%d", item.ID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(150, 8, item.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(87, 8, fmt.Sprintf("%d", item.PointsConversionRate), "1", 1, "L", false, 0, "")
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export member categories pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("member-categories-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
