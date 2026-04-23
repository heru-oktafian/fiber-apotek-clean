package unit

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type ConversionService struct {
	Units ports.UnitRepository
	IDs   ports.IDGenerator
}

func (s ConversionService) List(ctx context.Context, branchID string, req domain.ConversionListRequest) (domain.ConversionListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Units.ListConversions(ctx, branchID, req)
	if err != nil {
		return domain.ConversionListResult{}, apperror.New(http.StatusInternalServerError, "Get unit conversions failed", err.Error())
	}
	return items, nil
}

func (s ConversionService) GetByID(ctx context.Context, branchID, id string) (domain.ConversionMaster, error) {
	item, err := s.Units.FindConversionByID(ctx, id, branchID)
	if err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Get unit conversion failed", "unit conversion not found")
	}
	return item, nil
}

func (s ConversionService) Create(ctx context.Context, branchID string, req domain.ConversionCreateRequest) (domain.ConversionMaster, error) {
	productID := strings.TrimSpace(req.ProductID)
	initID := strings.TrimSpace(req.InitID)
	finalID := strings.TrimSpace(req.FinalID)
	if productID == "" || initID == "" || finalID == "" {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "product_id, init_id, and final_id are required")
	}
	if req.ValueConv <= 0 {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "value_conv must be greater than zero")
	}
	if initID == finalID {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "init_id and final_id must be different")
	}
	if _, err := s.Units.FindProductByID(ctx, productID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "product not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, initID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "initial unit not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, finalID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "final unit not found")
	}
	if _, err := s.Units.FindConversion(ctx, productID, initID, finalID, branchID); err == nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusConflict, "Create unit conversion failed", "unit conversion already exists")
	}
	item := domain.ConversionMaster{ID: s.IDs.New("UNC"), ProductID: productID, InitID: initID, FinalID: finalID, ValueConv: req.ValueConv, BranchID: branchID}
	if err := s.Units.CreateConversion(ctx, item); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusInternalServerError, "Create unit conversion failed", err.Error())
	}
	return s.Units.FindConversionByID(ctx, item.ID, branchID)
}

func (s ConversionService) Update(ctx context.Context, branchID, id string, req domain.ConversionCreateRequest) (domain.ConversionMaster, error) {
	productID := strings.TrimSpace(req.ProductID)
	initID := strings.TrimSpace(req.InitID)
	finalID := strings.TrimSpace(req.FinalID)
	if productID == "" || initID == "" || finalID == "" {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "product_id, init_id, and final_id are required")
	}
	if req.ValueConv <= 0 {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "value_conv must be greater than zero")
	}
	if initID == finalID {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "init_id and final_id must be different")
	}
	if _, err := s.Units.FindConversionByID(ctx, id, branchID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "unit conversion not found")
	}
	if _, err := s.Units.FindProductByID(ctx, productID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "product not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, initID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "initial unit not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, finalID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "final unit not found")
	}
	item := domain.ConversionMaster{ID: id, ProductID: productID, InitID: initID, FinalID: finalID, ValueConv: req.ValueConv, BranchID: branchID}
	if err := s.Units.UpdateConversion(ctx, item); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusInternalServerError, "Update unit conversion failed", err.Error())
	}
	return s.Units.FindConversionByID(ctx, id, branchID)
}

func (s ConversionService) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Units.FindConversionByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete unit conversion failed", "unit conversion not found")
	}
	if err := s.Units.DeleteConversion(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete unit conversion failed", err.Error())
	}
	return nil
}

func (s ConversionService) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Units.ListConversions(ctx, branchID, domain.ConversionListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export unit conversions excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Unit Conversions")
	sheet := "Unit Conversions"
	f.SetCellValue(sheet, "A1", "DATA UNIT CONVERSIONS")
	headers := []string{"ID", "PRODUCT", "INIT UNIT", "FINAL UNIT", "VALUE CONV"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.InitName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.FinalName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.ValueConv)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export unit conversions excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("unit-conversions-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s ConversionService) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Units.ListConversions(ctx, branchID, domain.ConversionListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export unit conversions pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER UNIT CONVERSIONS")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER UNIT CONVERSIONS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "PRODUCT", "INIT UNIT", "FINAL UNIT", "VALUE"}
	widths := []float64{35, 90, 55, 55, 42}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items.Items {
		values := []string{item.ID, item.ProductName, item.InitName, item.FinalName, fmt.Sprintf("%d", item.ValueConv)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export unit conversions pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("unit-conversions-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
