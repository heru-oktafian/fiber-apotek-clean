package buyreturn

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/buyreturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type Service struct {
	Repo  ports.BuyReturnRepository
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
	return s.Repo.ListBuyReturns(ctx, branchID, req)
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (domain.Detail, error) {
	header, err := s.Repo.FindBuyReturnByID(ctx, branchID, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Get buy return failed", "buy return not found")
	}
	items, err := s.Repo.FindBuyReturnItems(ctx, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Get buy return failed", err.Error())
	}
	formatted := make([]domain.FormattedItem, 0, len(items))
	for _, item := range items {
		formatted = append(formatted, domain.FormattedItem{
			ID:          item.ID,
			BuyReturnID: item.BuyReturnID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitID:      item.UnitID,
			UnitName:    item.UnitName,
			Price:       item.Price,
			Qty:         item.Qty,
			SubTotal:    item.SubTotal,
			ExpiredDate: item.ExpiredDate.Format("2006-01-02"),
		})
	}
	return domain.Detail{ID: header.ID, PurchaseID: header.PurchaseID, ReturnDate: header.ReturnDate.Format("2006-01-02"), TotalReturn: header.TotalReturn, Payment: header.Payment, Items: formatted}, nil
}

func (s Service) ListPurchaseSources(ctx context.Context, branchID, search, month string) ([]domain.PurchaseComboItem, error) {
	return s.Repo.ListPurchaseReturnSources(ctx, branchID, search, month)
}

func (s Service) ListReturnableItems(ctx context.Context, purchaseID string) ([]domain.ReturnableItem, error) {
	if strings.TrimSpace(purchaseID) == "" {
		return nil, apperror.New(http.StatusBadRequest, "Get buy return items failed", "purchase_id is required")
	}
	return s.Repo.ListPurchaseReturnableItems(ctx, purchaseID)
}

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, domain.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Buy Returns")
	sheet := "Buy Returns"
	f.SetCellValue(sheet, "A1", fmt.Sprintf("RETUR PEMBELIAN %s", month))
	headers := []string{"ID", "PEMBELIAN ID", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	grandTotal := 0
	for i, item := range result.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.PurchaseID)
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
		return nil, "", apperror.New(http.StatusInternalServerError, "Export buy returns excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("buy-returns-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, domain.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("RETUR PEMBELIAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("RETUR PEMBELIAN %s", month), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{45, 110, 40, 40, 42}
	headers := []string{"ID", "PEMBELIAN ID", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	grandTotal := 0
	for _, item := range result.Items {
		values := []string{item.ID, item.PurchaseID, item.ReturnDate, item.Payment, fmt.Sprintf("%d", item.TotalReturn)}
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
		return nil, "", apperror.New(http.StatusInternalServerError, "Export buy returns pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("RETUR-PEMBELIAN-%s.pdf", time.Now().Format("2006-01-02-15:04:05")), nil
}

func (s Service) ExportItemsExcel(ctx context.Context, branchID, buyReturnID string) ([]byte, string, error) {
	if strings.TrimSpace(buyReturnID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export buy return items excel failed", "buy_return_id is required")
	}
	detail, err := s.GetByID(ctx, branchID, buyReturnID)
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Detail Retur Pembelian")
	sheet := "Detail Retur Pembelian"
	f.SetCellValue(sheet, "A1", "LAPORAN DETAIL RETUR PEMBELIAN")
	f.SetCellValue(sheet, "A2", "ID RETUR PEMBELIAN")
	f.SetCellValue(sheet, "B2", ": "+detail.ID)
	f.SetCellValue(sheet, "A3", "TANGGAL RETUR")
	f.SetCellValue(sheet, "B3", ": "+detail.ReturnDate)
	f.SetCellValue(sheet, "A4", "ID PEMBELIAN")
	f.SetCellValue(sheet, "B4", ": "+detail.PurchaseID)
	f.SetCellValue(sheet, "A5", "METODE PEMBAYARAN")
	f.SetCellValue(sheet, "B5", ": "+detail.Payment)
	headers := []string{"PRODUK", "KADALUARSA", "QTY", "HARGA", "SUB TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s7", col), h)
	}
	for i, item := range detail.Items {
		row := i + 8
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.ExpiredDate)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("%d %s", item.Qty, item.UnitName))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.SubTotal)
	}
	totalRow := len(detail.Items) + 8
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), detail.TotalReturn)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export buy return items excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-RETUR-PEMBELIAN-%s-%s.xlsx", buyReturnID, time.Now().Format("20060102150405")), nil
}

func (s Service) ExportItemsPDF(ctx context.Context, branchID, buyReturnID string) ([]byte, string, error) {
	if strings.TrimSpace(buyReturnID) == "" {
		return nil, "", apperror.New(http.StatusBadRequest, "Export buy return items pdf failed", "buy_return_id is required")
	}
	detail, err := s.GetByID(ctx, branchID, buyReturnID)
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("DETAIL RETUR PEMBELIAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("ID RETUR PEMBELIAN : %s", detail.ID), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(277, 8, fmt.Sprintf("TANGGAL RETUR : %s", detail.ReturnDate), "", 1, "C", false, 0, "")
	pdf.CellFormat(277, 8, fmt.Sprintf("ID PEMBELIAN : %s | METODE PEMBAYARAN : %s", detail.PurchaseID, detail.Payment), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{115, 40, 35, 40, 47}
	headers := []string{"PRODUK", "KADALUARSA", "QTY", "HARGA", "SUB TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	for _, item := range detail.Items {
		values := []string{item.ProductName, item.ExpiredDate, fmt.Sprintf("%d %s", item.Qty, item.UnitName), fmt.Sprintf("%d", item.Price), fmt.Sprintf("%d", item.SubTotal)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(230, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(47, 8, fmt.Sprintf("%d", detail.TotalReturn), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export buy return items pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-RETUR-PEMBELIAN-%s.pdf", time.Now().Format("2006-01-02-15:04:05")), nil
}

func (s Service) Create(ctx context.Context, branchID, userID string, req domain.CreateRequest) (domain.Detail, error) {
	purchaseID := strings.TrimSpace(req.BuyReturn.PurchaseID)
	if purchaseID == "" {
		return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "purchase_id is required")
	}
	payment := strings.TrimSpace(req.BuyReturn.Payment)
	if payment == "" {
		payment = "paid_by_cash"
	}
	returnDate := s.Clock.Now()
	if strings.TrimSpace(req.BuyReturn.ReturnDate) != "" {
		parsed, err := time.Parse("2006-01-02", req.BuyReturn.ReturnDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "return_date must be in YYYY-MM-DD format")
		}
		returnDate = parsed
	}
	_, err := s.Repo.FindPurchaseByID(ctx, branchID, purchaseID)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Create buy return failed", "purchase not found")
	}
	headerID := s.IDs.New("BRT")
	header := domain.BuyReturn{ID: headerID, PurchaseID: purchaseID, ReturnDate: returnDate, BranchID: branchID, UserID: userID, Payment: payment, CreatedAt: s.Clock.Now(), UpdatedAt: s.Clock.Now()}
	items := make([]domain.Item, 0, len(req.BuyReturnItems))
	total := 0
	for _, in := range req.BuyReturnItems {
		purchaseItem, err := s.Repo.FindPurchaseItemByPurchaseAndProduct(ctx, purchaseID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "product not found in source purchase")
		}
		returnedQty, err := s.Repo.SumBuyReturnedQty(ctx, purchaseID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		if returnedQty+in.Qty > purchaseItem.Qty {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "return qty exceeds purchased qty")
		}
		expiredDate, err := time.Parse("2006-01-02", in.ExpiredDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "expired_date must be in YYYY-MM-DD format")
		}
		prod, err := s.Repo.FindProductByID(ctx, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		actualQty := in.Qty
		if purchaseItem.UnitID != "" && prod.UnitID != "" && purchaseItem.UnitID != prod.UnitID {
			conv, err := s.Repo.FindConversion(ctx, in.ProductID, purchaseItem.UnitID, prod.UnitID, branchID)
			if err == nil && conv.Value > 0 {
				actualQty = in.Qty * conv.Value
			}
		}
		prod.Stock -= actualQty
		if prod.Stock < 0 {
			prod.Stock = 0
		}
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		subTotal := purchaseItem.Price * in.Qty
		total += subTotal
		items = append(items, domain.Item{ID: s.IDs.New("BRI"), BuyReturnID: headerID, ProductID: in.ProductID, Price: purchaseItem.Price, Qty: in.Qty, SubTotal: subTotal, ExpiredDate: expiredDate})
	}
	header.TotalReturn = total
	if err := s.Repo.CreateBuyReturn(ctx, header); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	if err := s.Repo.CreateBuyReturnItems(ctx, items); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	if err := s.Repo.CreateTransactionReport(ctx, header.ID, "buy_return", userID, branchID, total, payment, s.Clock.Now()); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	return s.GetByID(ctx, branchID, headerID)
}
