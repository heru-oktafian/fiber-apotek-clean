package salereturn

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/salereturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type Service struct {
	Repo  ports.SaleReturnRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req domain.ListRequest) (domain.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	return s.Repo.ListSaleReturns(ctx, branchID, req)
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (domain.Detail, error) {
	header, err := s.Repo.FindSaleReturnByID(ctx, branchID, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Get sale return failed", "sale return not found")
	}
	items, err := s.Repo.FindSaleReturnItems(ctx, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Get sale return failed", err.Error())
	}
	formatted := make([]domain.FormattedItem, 0, len(items))
	for _, item := range items {
		formatted = append(formatted, domain.FormattedItem{
			ID:           item.ID,
			SaleReturnID: item.SaleReturnID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			UnitID:       item.UnitID,
			UnitName:     item.UnitName,
			Price:        item.Price,
			Qty:          item.Qty,
			SubTotal:     item.SubTotal,
			ExpiredDate:  item.ExpiredDate.Format("2006-01-02"),
		})
	}
	return domain.Detail{ID: header.ID, SaleID: header.SaleID, ReturnDate: header.ReturnDate.Format("2006-01-02"), TotalReturn: header.TotalReturn, Payment: header.Payment, Items: formatted}, nil
}

func (s Service) ListSaleSources(ctx context.Context, branchID, search, month string) ([]domain.SaleComboItem, error) {
	return s.Repo.ListSaleReturnSources(ctx, branchID, search, month)
}

func (s Service) ListReturnableItems(ctx context.Context, saleID string) ([]domain.ReturnableItem, error) {
	if strings.TrimSpace(saleID) == "" {
		return nil, apperror.New(http.StatusBadRequest, "Get sale return items failed", "sale_id is required")
	}
	return s.Repo.ListSaleReturnableItems(ctx, saleID)
}

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, domain.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Sale Returns")
	sheet := "Sale Returns"
	f.SetCellValue(sheet, "A1", fmt.Sprintf("RETUR PENJUALAN %s", month))
	headers := []string{"ID", "PENJUALAN ID", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	grandTotal := 0
	for i, item := range result.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.SaleID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.ReturnDate)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Payment)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.TotalReturn)
		grandTotal += item.TotalReturn
	}
	totalRow := len(result.Items) + 4
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "GRAND TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), grandTotal)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale returns excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("sale-returns-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, domain.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("RETUR PENJUALAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("RETUR PENJUALAN %s", month), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{45, 110, 40, 40, 42}
	headers := []string{"ID", "PENJUALAN ID", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	grandTotal := 0
	for _, item := range result.Items {
		values := []string{item.ID, item.SaleID, item.ReturnDate, item.Payment, fmt.Sprintf("%d", item.TotalReturn)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		grandTotal += item.TotalReturn
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(235, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(42, 8, fmt.Sprintf("%d", grandTotal), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale returns pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("RETUR-PENJUALAN-%s.pdf", time.Now().Format("2006-01-02-15:04:05")), nil
}

func (s Service) ExportItemsExcel(ctx context.Context, branchID, saleReturnID string) ([]byte, string, error) {
	if strings.TrimSpace(saleReturnID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export sale return items excel failed", "sale_return_id is required")
	}
	detail, err := s.GetByID(ctx, branchID, saleReturnID)
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Detail Retur Penjualan")
	sheet := "Detail Retur Penjualan"
	f.SetCellValue(sheet, "A1", "LAPORAN DETAIL RETUR PENJUALAN")
	f.SetCellValue(sheet, "A2", "ID RETUR PENJUALAN")
	f.SetCellValue(sheet, "B2", ": "+detail.ID)
	f.SetCellValue(sheet, "A3", "TANGGAL RETUR")
	f.SetCellValue(sheet, "B3", ": "+detail.ReturnDate)
	f.SetCellValue(sheet, "A4", "ID PENJUALAN")
	f.SetCellValue(sheet, "B4", ": "+detail.SaleID)
	f.SetCellValue(sheet, "A5", "METODE PEMBAYARAN")
	f.SetCellValue(sheet, "B5", ": "+detail.Payment)
	headers := []string{"PRODUK", "KADALUARSA", "HARGA", "QTY", "SUB TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s7", col), h)
	}
	for i, item := range detail.Items {
		row := i + 8
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.ExpiredDate)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Price)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("%d %s", item.Qty, item.UnitName))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.SubTotal)
	}
	totalRow := len(detail.Items) + 8
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), detail.TotalReturn)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale return items excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-RETUR-PENJUALAN-%s-%s.xlsx", saleReturnID, time.Now().Format("20060102150405")), nil
}

func (s Service) ExportItemsPDF(ctx context.Context, branchID, saleReturnID string) ([]byte, string, error) {
	if strings.TrimSpace(saleReturnID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export sale return items pdf failed", "sale_return_id is required")
	}
	detail, err := s.GetByID(ctx, branchID, saleReturnID)
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("DETAIL RETUR PENJUALAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("ID RETUR PENJUALAN : %s", detail.ID), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(277, 8, fmt.Sprintf("TANGGAL RETUR : %s", detail.ReturnDate), "", 1, "C", false, 0, "")
	pdf.CellFormat(277, 8, fmt.Sprintf("ID PENJUALAN : %s | METODE PEMBAYARAN : %s", detail.SaleID, detail.Payment), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{105, 40, 40, 35, 57}
	headers := []string{"PRODUK", "KADALUARSA", "HARGA", "QTY", "SUB TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	for _, item := range detail.Items {
		values := []string{item.ProductName, item.ExpiredDate, fmt.Sprintf("%d", item.Price), fmt.Sprintf("%d %s", item.Qty, item.UnitName), fmt.Sprintf("%d", item.SubTotal)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(220, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(57, 8, fmt.Sprintf("%d", detail.TotalReturn), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale return items pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-RETUR-PENJUALAN-%s-%s.pdf", saleReturnID, time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) Create(ctx context.Context, branchID, userID string, req domain.CreateRequest) (domain.Detail, error) {
	saleID := strings.TrimSpace(req.SaleReturn.SaleID)
	if saleID == "" {
		return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "sale_id is required")
	}
	payment := strings.TrimSpace(req.SaleReturn.Payment)
	if payment == "" {
		payment = "paid_by_cash"
	}
	returnDate := s.Clock.Now()
	if strings.TrimSpace(req.SaleReturn.ReturnDate) != "" {
		parsed, err := time.Parse("2006-01-02", req.SaleReturn.ReturnDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "return_date must be in YYYY-MM-DD format")
		}
		returnDate = parsed
	}
	_, err := s.Repo.FindSaleByID(ctx, branchID, saleID)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Create sale return failed", "sale not found")
	}
	headerID := s.IDs.New("SRT")
	header := domain.SaleReturn{ID: headerID, SaleID: saleID, ReturnDate: returnDate, BranchID: branchID, UserID: userID, Payment: payment, CreatedAt: s.Clock.Now(), UpdatedAt: s.Clock.Now()}
	items := make([]domain.Item, 0, len(req.SaleReturnItems))
	total := 0
	for _, in := range req.SaleReturnItems {
		saleItem, err := s.Repo.FindSaleItemBySaleAndProduct(ctx, saleID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "product not found in source sale")
		}
		returnedQty, err := s.Repo.SumSaleReturnedQty(ctx, saleID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		if returnedQty+in.Qty > saleItem.Qty {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "return qty exceeds sold qty")
		}
		expiredDate, err := time.Parse("2006-01-02", in.ExpiredDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "expired_date must be in YYYY-MM-DD format")
		}
		prod, err := s.Repo.FindProductByID(ctx, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		prod.Stock += in.Qty
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		subTotal := saleItem.Price * in.Qty
		total += subTotal
		items = append(items, domain.Item{ID: s.IDs.New("SRI"), SaleReturnID: headerID, ProductID: in.ProductID, Price: saleItem.Price, Qty: in.Qty, SubTotal: subTotal, ExpiredDate: expiredDate})
	}
	header.TotalReturn = total
	if err := s.Repo.CreateSaleReturn(ctx, header); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	if err := s.Repo.CreateSaleReturnItems(ctx, items); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	if err := s.Repo.CreateTransactionReport(ctx, header.ID, "sale_return", userID, branchID, total, payment, s.Clock.Now()); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	return s.GetByID(ctx, branchID, headerID)
}
