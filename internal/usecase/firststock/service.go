package firststock

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/firststock"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
	"gorm.io/gorm"
)

type Service struct {
	Repo  ports.FirstStockRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req firststock.ListRequest) (firststock.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Month == "" {
		req.Month = s.Clock.Now().Format("2006-01")
	}
	return s.Repo.ListFirstStocks(ctx, branchID, req)
}

func (s Service) Create(ctx context.Context, branchID, userID string, req firststock.CreateRequest) (firststock.FirstStock, error) {
	now := s.Clock.Now()
	date := now
	if strings.TrimSpace(req.FirstStockDate) != "" {
		parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.FirstStockDate))
		if err != nil {
			return firststock.FirstStock{}, apperror.New(http.StatusBadRequest, "Create first stock failed", "invalid date format, use YYYY-MM-DD")
		}
		date = parsedDate
	}
	item := firststock.FirstStock{ID: s.IDs.New("FST"), Description: strings.TrimSpace(req.Description), FirstStockDate: date, BranchID: branchID, UserID: userID, TotalFirstStock: 0, Payment: common.PaymentStatus("nocost"), CreatedAt: now, UpdatedAt: now}
	if err := s.Repo.CreateFirstStock(ctx, item); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Create first stock failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "first_stock", item.UserID, item.BranchID, item.TotalFirstStock, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Create first stock failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, req firststock.UpdateRequest) (firststock.FirstStock, error) {
	item, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusNotFound, "Update first stock failed", "first stock not found")
	}
	if strings.TrimSpace(req.FirstStockDate) != "" {
		parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.FirstStockDate))
		if err != nil {
			return firststock.FirstStock{}, apperror.New(http.StatusBadRequest, "Update first stock failed", "invalid date format, use YYYY-MM-DD")
		}
		item.FirstStockDate = parsedDate
	}
	if strings.TrimSpace(req.Description) != "" {
		item.Description = strings.TrimSpace(req.Description)
	}
	if strings.TrimSpace(req.Payment) != "" {
		item.Payment = common.PaymentStatus(strings.TrimSpace(req.Payment))
	}
	items, err := s.Repo.FindFirstStockItems(ctx, item.ID)
	if err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	item.TotalFirstStock = 0
	for _, line := range items {
		item.TotalFirstStock += line.SubTotal
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateFirstStock(ctx, item); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "first_stock", item.UserID, item.BranchID, item.TotalFirstStock, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	return item, nil
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete first stock failed", "first stock not found")
	}
	items, err := s.Repo.FindFirstStockItems(ctx, id)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	for _, item := range items {
		prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
		prod.Stock -= item.Qty
		if prod.Stock < 0 {
			prod.Stock = 0
		}
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
	}
	for _, item := range items {
		if err := s.Repo.DeleteFirstStockItem(ctx, item.ID); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
	}
	if err := s.Repo.DeleteTransactionReport(ctx, header.ID, "first_stock"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	if err := s.Repo.DeleteFirstStock(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	return nil
}

func (s Service) GetDetail(ctx context.Context, branchID, id string) (firststock.Detail, error) {
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return firststock.Detail{}, apperror.New(http.StatusNotFound, "Get first stock detail failed", "first stock not found")
	}
	items, err := s.Repo.FindFirstStockItems(ctx, id)
	if err != nil {
		return firststock.Detail{}, apperror.New(http.StatusInternalServerError, "Get first stock detail failed", err.Error())
	}
	return firststock.Detail{ID: header.ID, Description: header.Description, FirstStockDate: header.FirstStockDate.Format("02-01-2006"), TotalFirstStock: header.TotalFirstStock, Payment: string(header.Payment), Items: items}, nil
}

func (s Service) ListItems(ctx context.Context, firstStockID string) ([]firststock.Item, error) {
	items, err := s.Repo.FindFirstStockItems(ctx, firstStockID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "List first stock items failed", err.Error())
	}
	return items, nil
}

func (s Service) CreateItem(ctx context.Context, branchID string, req firststock.CreateItemRequest) (firststock.Item, error) {
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusBadRequest, "Create first stock item failed", "invalid expired_date format, use YYYY-MM-DD")
	}
	_, err = s.Repo.FindFirstStockByID(ctx, branchID, req.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "first stock not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "product not found")
	}
	unitItem, err := s.Repo.FindUnit(ctx, req.UnitID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "unit not found")
	}
	conversion := 1
	if req.UnitID != prod.UnitID {
		conv, err := s.Repo.FindConversion(ctx, req.ProductID, req.UnitID, prod.UnitID, branchID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
		}
		if err == nil {
			conversion = conv.Value
		}
	}
	actualQty := req.Qty * conversion
	item := firststock.Item{ID: s.IDs.New("FSI"), FirstStockID: req.FirstStockID, ProductID: req.ProductID, ProductName: prod.Name, UnitID: req.UnitID, UnitName: unitItem.Name, Price: prod.PurchasePrice, Qty: req.Qty, SubTotal: prod.PurchasePrice * actualQty, ExpiredDate: expiredDate}
	if err := s.Repo.CreateFirstStockItem(ctx, item); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	prod.Stock += actualQty
	if expiredDate.Before(prod.ExpiredDate) || prod.ExpiredDate.IsZero() {
		prod.ExpiredDate = expiredDate
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, req.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, req.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return item, nil
}

func (s Service) UpdateItem(ctx context.Context, branchID, id string, req firststock.UpdateItemRequest) (firststock.Item, error) {
	existing, err := s.Repo.FindFirstStockItemByID(ctx, id)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "item not found")
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusBadRequest, "Update first stock item failed", "invalid expired_date format, use YYYY-MM-DD")
	}
	oldProduct, err := s.Repo.FindProductByID(ctx, existing.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	oldProduct.Stock -= existing.Qty
	if oldProduct.Stock < 0 {
		oldProduct.Stock = 0
	}
	if err := s.Repo.UpdateProduct(ctx, oldProduct); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "product not found")
	}
	unitItem, err := s.Repo.FindUnit(ctx, req.UnitID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "unit not found")
	}
	conversion := 1
	if req.UnitID != prod.UnitID {
		conv, err := s.Repo.FindConversion(ctx, req.ProductID, req.UnitID, prod.UnitID, branchID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
		}
		if err == nil {
			conversion = conv.Value
		}
	}
	actualQty := req.Qty * conversion
	existing.ProductID = req.ProductID
	existing.ProductName = prod.Name
	existing.UnitID = req.UnitID
	existing.UnitName = unitItem.Name
	existing.Price = prod.PurchasePrice
	existing.Qty = req.Qty
	existing.SubTotal = prod.PurchasePrice * actualQty
	existing.ExpiredDate = expiredDate
	if err := s.Repo.UpdateFirstStockItem(ctx, existing); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	prod.Stock += actualQty
	if expiredDate.Before(prod.ExpiredDate) || prod.ExpiredDate.IsZero() {
		prod.ExpiredDate = expiredDate
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, existing.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, existing.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return existing, nil
}

func (s Service) DeleteItem(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindFirstStockItemByID(ctx, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete first stock item failed", "item not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	prod.Stock -= item.Qty
	if prod.Stock < 0 {
		prod.Stock = 0
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	if err := s.Repo.DeleteFirstStockItem(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, item.FirstStockID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, item.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return nil
}

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, firststock.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("First Stocks")
	sheet := "First Stocks"
	f.SetCellValue(sheet, "A1", fmt.Sprintf("STOK AWAL %s", month))
	headers := []string{"ID", "KETERANGAN", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	grandTotal := 0
	for i, item := range result.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Description)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.FirstStockDate)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Payment)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.TotalFirstStock)
		grandTotal += item.TotalFirstStock
	}
	totalRow := len(result.Items) + 4
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "GRAND TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), grandTotal)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export first stocks excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("first-stocks-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, firststock.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("STOK AWAL")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("STOK AWAL %s", month), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{45, 110, 40, 40, 42}
	headers := []string{"ID", "KETERANGAN", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	grandTotal := 0
	for _, item := range result.Items {
		values := []string{item.ID, item.Description, item.FirstStockDate, item.Payment, fmt.Sprintf("%d", item.TotalFirstStock)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		grandTotal += item.TotalFirstStock
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(235, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(42, 8, fmt.Sprintf("%d", grandTotal), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export first stocks pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("STOK-AWAL-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportItemsExcel(ctx context.Context, branchID, firstStockID string) ([]byte, string, error) {
	if strings.TrimSpace(firstStockID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export first stock items excel failed", "first_stock_id is required")
	}
	header, err := s.GetDetail(ctx, branchID, firstStockID)
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Detail Stok Awal")
	sheet := "Detail Stok Awal"
	f.SetCellValue(sheet, "A1", "LAPORAN DETAIL STOK AWAL")
	f.SetCellValue(sheet, "A2", "ID STOK AWAL")
	f.SetCellValue(sheet, "B2", ": "+header.ID)
	f.SetCellValue(sheet, "A3", "TANGGAL")
	f.SetCellValue(sheet, "B3", ": "+header.FirstStockDate)
	f.SetCellValue(sheet, "A4", "METODE PEMBAYARAN")
	f.SetCellValue(sheet, "B4", ": "+header.Payment)
	f.SetCellValue(sheet, "A5", "DESKRIPSI")
	f.SetCellValue(sheet, "B5", ": "+header.Description)
	headers := []string{"PRODUK", "QTY", "HARGA", "SUB TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s7", col), h)
	}
	for i, item := range header.Items {
		row := i + 8
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%d %s", item.Qty, item.UnitName))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Price)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.SubTotal)
	}
	totalRow := len(header.Items) + 8
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", totalRow), header.TotalFirstStock)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export first stock items excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-STOK-AWAL-%s-%s.xlsx", firstStockID, time.Now().Format("20060102150405")), nil
}

func (s Service) ExportItemsPDF(ctx context.Context, branchID, firstStockID string) ([]byte, string, error) {
	if strings.TrimSpace(firstStockID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export first stock items pdf failed", "first_stock_id is required")
	}
	header, err := s.GetDetail(ctx, branchID, firstStockID)
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("DETAIL STOK AWAL")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("STOK AWAL : %s", header.ID), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(277, 8, fmt.Sprintf("TANGGAL : %s | METODE PEMBAYARAN : %s", header.FirstStockDate, header.Payment), "", 1, "C", false, 0, "")
	pdf.CellFormat(277, 8, fmt.Sprintf("DESKRIPSI : %s", header.Description), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{130, 45, 50, 52}
	headers := []string{"PRODUK", "QTY", "HARGA", "SUB TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	for _, item := range header.Items {
		values := []string{item.ProductName, fmt.Sprintf("%d %s", item.Qty, item.UnitName), fmt.Sprintf("%d", item.Price), fmt.Sprintf("%d", item.SubTotal)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(225, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(52, 8, fmt.Sprintf("%d", header.TotalFirstStock), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export first stock items pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-STOK-AWAL-%s.pdf", time.Now().Format("2006-01-02-15:04:05")), nil
}
