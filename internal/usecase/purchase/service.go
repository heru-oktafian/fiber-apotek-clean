package purchase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type Service struct {
	Repo  ports.PurchaseRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req purchase.ListRequest) (purchase.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Repo.ListPurchases(ctx, branchID, req)
	if err != nil {
		return purchase.ListResult{}, apperror.New(http.StatusInternalServerError, "Get purchases failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (purchase.Detail, error) {
	item, err := s.Repo.FindPurchaseDetail(ctx, branchID, id)
	if err != nil {
		return purchase.Detail{}, apperror.New(http.StatusNotFound, "Get purchase failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, req purchase.UpdateRequest) (purchase.Purchase, error) {
	item, err := s.Repo.FindPurchaseByID(ctx, branchID, id)
	if err != nil {
		return purchase.Purchase{}, apperror.New(http.StatusNotFound, "Update purchase failed", "purchase not found")
	}
	if req.SupplierID != "" {
		item.SupplierID = req.SupplierID
	}
	if req.PurchaseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			return purchase.Purchase{}, apperror.New(http.StatusBadRequest, "Update purchase failed", "invalid purchase_date format. use YYYY-MM-DD")
		}
		item.PurchaseDate = parsedDate
	}
	if req.Payment != "" {
		item.Payment = common.PaymentStatus(req.Payment)
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdatePurchaseHeader(ctx, item); err != nil {
		return purchase.Purchase{}, apperror.New(http.StatusInternalServerError, "Update purchase failed", err.Error())
	}
	return s.Repo.FindPurchaseByID(ctx, branchID, id)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Repo.FindPurchaseByID(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusNotFound, "Delete purchase failed", "purchase not found")
	}
	if err := s.Repo.DeletePurchaseHeader(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete purchase failed", err.Error())
	}
	return nil
}

func (s Service) ListItems(ctx context.Context, branchID, purchaseID string) ([]purchase.Item, error) {
	if _, err := s.Repo.FindPurchaseByID(ctx, branchID, purchaseID); err != nil {
		return nil, apperror.New(http.StatusNotFound, "Get purchase items failed", "purchase not found")
	}
	items, err := s.Repo.FindPurchaseItems(ctx, purchaseID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get purchase items failed", err.Error())
	}
	return items, nil
}

func (s Service) CreateItem(ctx context.Context, branchID string, req purchase.CreateItemRequest) (purchase.Item, error) {
	header, err := s.Repo.FindPurchaseByID(ctx, branchID, req.PurchaseID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusNotFound, "Create purchase item failed", "purchase not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusNotFound, "Create purchase item failed", "product not found")
	}
	if req.UnitID == "" {
		req.UnitID = prod.UnitID
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusBadRequest, "Create purchase item failed", "invalid expired_date format. use YYYY-MM-DD")
	}
	items, err := s.Repo.FindPurchaseItems(ctx, req.PurchaseID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
	}
	for _, existing := range items {
		if existing.ProductID == req.ProductID {
			existing.Qty += req.Qty
			existing.SubTotal = existing.Qty * existing.Price
			if err := s.Repo.UpdatePurchaseItem(ctx, existing); err != nil {
				return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
			}
			prod.Stock += req.Qty
			if req.Price > prod.PurchasePrice {
				prod.PurchasePrice = req.Price
			}
			prod.ExpiredDate = expiredDate
			if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
				return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
			}
			header.TotalPurchase += req.Price * req.Qty
			header.UpdatedAt = s.Clock.Now()
			if err := s.Repo.UpdatePurchaseHeader(ctx, header); err != nil {
				return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
			}
			return s.Repo.FindPurchaseItemByID(ctx, existing.ID)
		}
	}
	item := purchase.Item{ID: s.IDs.New("PIT"), PurchaseID: req.PurchaseID, ProductID: req.ProductID, UnitID: req.UnitID, Price: req.Price, Qty: req.Qty, SubTotal: req.Price * req.Qty, ExpiredDate: expiredDate}
	if err := s.Repo.CreatePurchaseItem(ctx, item); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
	}
	prod.Stock += req.Qty
	if req.Price > prod.PurchasePrice {
		prod.PurchasePrice = req.Price
	}
	prod.ExpiredDate = expiredDate
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
	}
	header.TotalPurchase += item.SubTotal
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdatePurchaseHeader(ctx, header); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Create purchase item failed", err.Error())
	}
	return s.Repo.FindPurchaseItemByID(ctx, item.ID)
}

func (s Service) UpdateItem(ctx context.Context, branchID, id string, req purchase.UpdateItemRequest) (purchase.Item, error) {
	item, err := s.Repo.FindPurchaseItemByID(ctx, id)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusNotFound, "Update purchase item failed", "item not found")
	}
	header, err := s.Repo.FindPurchaseByID(ctx, branchID, item.PurchaseID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusNotFound, "Update purchase item failed", "purchase not found")
	}
	oldSubtotal := item.SubTotal
	oldProductID := item.ProductID
	oldQty := item.Qty
	prodOld, err := s.Repo.FindProductByID(ctx, oldProductID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Update purchase item failed", err.Error())
	}
	prodOld.Stock -= oldQty
	if err := s.Repo.UpdateProduct(ctx, prodOld); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Update purchase item failed", err.Error())
	}
	prodNew, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusNotFound, "Update purchase item failed", "product not found")
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return purchase.Item{}, apperror.New(http.StatusBadRequest, "Update purchase item failed", "invalid expired_date format. use YYYY-MM-DD")
	}
	prodNew.Stock += req.Qty
	if req.Price > prodNew.PurchasePrice {
		prodNew.PurchasePrice = req.Price
	}
	prodNew.ExpiredDate = expiredDate
	if err := s.Repo.UpdateProduct(ctx, prodNew); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Update purchase item failed", err.Error())
	}
	item.ProductID = req.ProductID
	item.UnitID = req.UnitID
	item.Price = req.Price
	item.Qty = req.Qty
	item.SubTotal = req.Price * req.Qty
	item.ExpiredDate = expiredDate
	if err := s.Repo.UpdatePurchaseItem(ctx, item); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Update purchase item failed", err.Error())
	}
	header.TotalPurchase += item.SubTotal - oldSubtotal
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdatePurchaseHeader(ctx, header); err != nil {
		return purchase.Item{}, apperror.New(http.StatusInternalServerError, "Update purchase item failed", err.Error())
	}
	return s.Repo.FindPurchaseItemByID(ctx, id)
}

func (s Service) DeleteItem(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindPurchaseItemByID(ctx, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete purchase item failed", "item not found")
	}
	header, err := s.Repo.FindPurchaseByID(ctx, branchID, item.PurchaseID)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete purchase item failed", "purchase not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete purchase item failed", err.Error())
	}
	prod.Stock -= item.Qty
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete purchase item failed", err.Error())
	}
	if err := s.Repo.DeletePurchaseItem(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete purchase item failed", err.Error())
	}
	header.TotalPurchase -= item.SubTotal
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdatePurchaseHeader(ctx, header); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete purchase item failed", err.Error())
	}
	return nil
}

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	items, err := s.Repo.ListPurchases(ctx, branchID, purchase.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchases excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Purchases")
	sheet := "Purchases"
	f.SetCellValue(sheet, "A1", "DATA PEMBELIAN")
	headers := []string{"ID", "SUPPLIER", "DATE", "TOTAL", "PAYMENT"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.SupplierName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.PurchaseDate.Format("02/01/2006"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.TotalPurchase)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Payment)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchases excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("purchases-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	items, err := s.Repo.ListPurchases(ctx, branchID, purchase.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchases pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("PEMBELIAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "PEMBELIAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{45, 80, 40, 50, 62}
	headers := []string{"ID", "SUPPLIER", "DATE", "TOTAL", "PAYMENT"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	for _, item := range items.Items {
		values := []string{item.ID, item.SupplierName, item.PurchaseDate.Format("02/01/2006"), fmt.Sprintf("%d", item.TotalPurchase), string(item.Payment)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchases pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("PEMBELIAN-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportItemsExcel(ctx context.Context, branchID, purchaseID string) ([]byte, string, error) {
	items, err := s.ListItems(ctx, branchID, purchaseID)
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Purchase Items")
	sheet := "Purchase Items"
	f.SetCellValue(sheet, "A1", "DETAIL PEMBELIAN")
	headers := []string{"ID", "PRODUCT", "UNIT", "PRICE", "QTY", "SUB TOTAL", "EXPIRED DATE"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.UnitName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Qty)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.SubTotal)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), item.ExpiredDate.Format("02/01/2006"))
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchase items excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-PEMBELIAN-%s-%s.xlsx", purchaseID, time.Now().Format("20060102150405")), nil
}

func (s Service) ExportItemsPDF(ctx context.Context, branchID, purchaseID string) ([]byte, string, error) {
	items, err := s.ListItems(ctx, branchID, purchaseID)
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("DETAIL PEMBELIAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "DETAIL PEMBELIAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	widths := []float64{25, 70, 30, 30, 20, 35, 67}
	headers := []string{"ID", "PRODUCT", "UNIT", "PRICE", "QTY", "SUB TOTAL", "EXPIRED DATE"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items {
		values := []string{item.ID, item.ProductName, item.UnitName, fmt.Sprintf("%d", item.Price), fmt.Sprintf("%d", item.Qty), fmt.Sprintf("%d", item.SubTotal), item.ExpiredDate.Format("02/01/2006")}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export purchase items pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-PEMBELIAN-%s.pdf", purchaseID), nil
}

func (s Service) CreateTransaction(ctx context.Context, branchID, userID string, req purchase.CreatePurchaseRequest) (purchase.Purchase, []purchase.Item, error) {
	now := s.Clock.Now()
	purchaseDate := now
	if req.Purchase.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.Purchase.PurchaseDate)
		if err != nil {
			return purchase.Purchase{}, nil, apperror.New(http.StatusBadRequest, "Invalid purchase_date format. Please use `YYYY-MM-DD`.", err)
		}
		purchaseDate = parsed
	}
	payment := common.PaymentStatus(req.Purchase.Payment)
	if payment == "" {
		payment = common.PaymentCash
	}
	p := purchase.Purchase{
		ID:           s.IDs.New("PUR"),
		SupplierID:   req.Purchase.SupplierID,
		PurchaseDate: purchaseDate,
		BranchID:     branchID,
		UserID:       userID,
		Payment:      payment,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	var items []purchase.Item
	if err := s.Repo.WithinTransaction(ctx, func(repo ports.PurchaseTxRepository) error {
		var total int
		for _, input := range req.PurchaseItems {
			prod, err := repo.FindProduct(ctx, input.ProductID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Product with ID %s not found", input.ProductID), err)
			}
			_, err = repo.FindUnit(ctx, input.UnitID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Unit with ID %s not found", input.UnitID), err)
			}
			conversionValue := 1
			if input.UnitID != prod.UnitID {
				conv, err := repo.FindConversion(ctx, input.ProductID, input.UnitID, prod.UnitID, branchID)
				if err == nil && conv.Value > 0 {
					conversionValue = conv.Value
				}
			}
			expiredDate, err := time.Parse("2006-01-02", input.ExpiredDate)
			if err != nil {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("Invalid expired_date format for product %s. Please use `YYYY-MM-DD`.", input.ProductID), err)
			}
			itemPrice := input.Price * conversionValue
			itemSubtotal := itemPrice * input.Qty
			items = append(items, purchase.Item{
				ID:          s.IDs.New("PIT"),
				PurchaseID:  p.ID,
				ProductID:   input.ProductID,
				UnitID:      input.UnitID,
				Price:       itemPrice,
				Qty:         input.Qty,
				SubTotal:    itemSubtotal,
				ExpiredDate: expiredDate,
			})
			actualQtyToAdd := input.Qty * conversionValue
			prod.Stock += actualQtyToAdd
			if prod.ExpiredDate.IsZero() || expiredDate.Before(prod.ExpiredDate) {
				prod.ExpiredDate = expiredDate
			}
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update product details (stock/expired_date) for product %s", prod.Name), err)
			}
			total += itemSubtotal
		}
		p.TotalPurchase = total
		if err := repo.CreatePurchase(ctx, p); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create purchase", err)
		}
		if err := repo.CreatePurchaseItems(ctx, items); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create purchase items", err)
		}
		if err := repo.CreateTransactionReport(ctx, s.IDs.New("TRX"), "purchase", userID, branchID, p.TotalPurchase, string(p.Payment), now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create transaction report", err)
		}
		return nil
	}); err != nil {
		return purchase.Purchase{}, nil, err
	}
	return p, items, nil
}
