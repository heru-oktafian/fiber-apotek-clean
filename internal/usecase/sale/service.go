package sale

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
)

type Service struct {
	Repo  ports.SaleRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req sale.ListRequest) (sale.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Repo.ListSales(ctx, branchID, req)
	if err != nil {
		return sale.ListResult{}, apperror.New(http.StatusInternalServerError, "Get sales failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (sale.Detail, error) {
	item, err := s.Repo.FindSaleDetail(ctx, branchID, id)
	if err != nil {
		return sale.Detail{}, apperror.New(http.StatusNotFound, "Get sale failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id, defaultMember string, req sale.UpdateRequest) (sale.Sale, error) {
	item, err := s.Repo.FindSaleByID(ctx, branchID, id)
	if err != nil {
		return sale.Sale{}, apperror.New(http.StatusNotFound, "Update sale failed", "sale not found")
	}
	oldTotal := item.TotalSale
	oldProfit := item.ProfitEstimate
	if req.MemberID != nil {
		memberID := *req.MemberID
		if memberID == "" {
			memberID = defaultMember
		}
		if _, err := s.Repo.FindMemberByID(ctx, memberID); err != nil {
			memberID = defaultMember
		}
		item.MemberID = memberID
	}
	if req.Payment != "" {
		item.Payment = common.PaymentStatus(req.Payment)
	}
	if req.Discount != nil {
		item.Discount = *req.Discount
	}
	items, err := s.Repo.FindSaleItems(ctx, id)
	if err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	var total int
	var profit int
	for _, line := range items {
		total += line.SubTotal
		profit += line.SubTotal - (line.Price * line.Qty)
	}
	item.TotalSale = total - item.Discount
	item.ProfitEstimate = profit
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateSaleHeader(ctx, item); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, item.ID, item.TotalSale, string(item.Payment), item.UpdatedAt); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.SaleDate, item.UserID, item.BranchID, item.TotalSale-oldTotal, item.ProfitEstimate-oldProfit, item.UpdatedAt); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	return s.Repo.FindSaleByID(ctx, branchID, id)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindSaleByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete sale failed", "sale not found")
	}
	items, err := s.Repo.FindSaleItems(ctx, id)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	for _, line := range items {
		prod, err := s.Repo.FindProductByID(ctx, line.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
		}
		prod.Stock += line.Qty
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
		}
	}
	now := s.Clock.Now()
	if err := s.Repo.DeleteSaleItems(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.DeleteTransactionReport(ctx, id, "sale"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.DeleteSaleHeader(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.SaleDate, item.UserID, item.BranchID, -item.TotalSale, -item.ProfitEstimate, now); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	return nil
}

func (s Service) ListItems(ctx context.Context, branchID, saleID string) ([]sale.Item, error) {
	if _, err := s.Repo.FindSaleByID(ctx, branchID, saleID); err != nil {
		return nil, apperror.New(http.StatusNotFound, "Get sale items failed", "sale not found")
	}
	items, err := s.Repo.FindSaleItems(ctx, saleID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get sale items failed", err.Error())
	}
	return items, nil
}

func (s Service) CreateItem(ctx context.Context, branchID string, req sale.CreateItemRequest) (sale.Item, error) {
	header, err := s.Repo.FindSaleByID(ctx, branchID, req.SaleID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusNotFound, "Create sale item failed", "sale not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusNotFound, "Create sale item failed", "product not found")
	}
	if prod.Stock < req.Qty {
		return sale.Item{}, apperror.New(http.StatusBadRequest, "Create sale item failed", fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", prod.Name, prod.Stock, req.Qty))
	}
	items, err := s.Repo.FindSaleItems(ctx, req.SaleID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	oldTotal := header.TotalSale
	oldProfit := header.ProfitEstimate
	for _, existing := range items {
		if existing.ProductID == req.ProductID {
			existing.Qty += req.Qty
			existing.Price = prod.SalesPrice
			existing.SubTotal = existing.Qty * existing.Price
			if err := s.Repo.UpdateSaleItem(ctx, existing); err != nil {
				return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
			}
			prod.Stock -= req.Qty
			if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
				return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
			}
			header.TotalSale += prod.SalesPrice * req.Qty
			header.ProfitEstimate += (prod.SalesPrice - prod.PurchasePrice) * req.Qty
			header.UpdatedAt = s.Clock.Now()
			if err := s.Repo.UpdateSaleHeader(ctx, header); err != nil {
				return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
			}
			if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalSale, string(header.Payment), header.UpdatedAt); err != nil {
				return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
			}
			if err := s.Repo.AdjustDailyProfit(ctx, header.SaleDate, header.UserID, header.BranchID, header.TotalSale-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
				return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
			}
			return s.Repo.FindSaleItemByID(ctx, existing.ID)
		}
	}
	item := sale.Item{ID: s.IDs.New("SIT"), SaleID: req.SaleID, ProductID: req.ProductID, Price: prod.SalesPrice, Qty: req.Qty, SubTotal: prod.SalesPrice * req.Qty}
	if err := s.Repo.CreateSaleItem(ctx, item); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	prod.Stock -= req.Qty
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	header.TotalSale += item.SubTotal
	header.ProfitEstimate += (item.Price - prod.PurchasePrice) * item.Qty
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateSaleHeader(ctx, header); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalSale, string(header.Payment), header.UpdatedAt); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.SaleDate, header.UserID, header.BranchID, header.TotalSale-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Create sale item failed", err.Error())
	}
	return s.Repo.FindSaleItemByID(ctx, item.ID)
}

func (s Service) UpdateItem(ctx context.Context, branchID, id string, req sale.UpdateItemRequest) (sale.Item, error) {
	item, err := s.Repo.FindSaleItemByID(ctx, id)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusNotFound, "Update sale item failed", "item not found")
	}
	header, err := s.Repo.FindSaleByID(ctx, branchID, item.SaleID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusNotFound, "Update sale item failed", "sale not found")
	}
	oldTotal := header.TotalSale
	oldProfit := header.ProfitEstimate
	oldProduct, err := s.Repo.FindProductByID(ctx, item.ProductID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	oldProduct.Stock += item.Qty
	if err := s.Repo.UpdateProduct(ctx, oldProduct); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	newProduct, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusNotFound, "Update sale item failed", "product not found")
	}
	if newProduct.Stock < req.Qty {
		return sale.Item{}, apperror.New(http.StatusBadRequest, "Update sale item failed", fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", newProduct.Name, newProduct.Stock, req.Qty))
	}
	newProduct.Stock -= req.Qty
	if err := s.Repo.UpdateProduct(ctx, newProduct); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	item.ProductID = req.ProductID
	item.Price = newProduct.SalesPrice
	item.Qty = req.Qty
	item.SubTotal = newProduct.SalesPrice * req.Qty
	if err := s.Repo.UpdateSaleItem(ctx, item); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	items, err := s.Repo.FindSaleItems(ctx, item.SaleID)
	if err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	var total int
	var profit int
	for _, line := range items {
		total += line.SubTotal
		profit += line.SubTotal - (line.Price * line.Qty)
	}
	header.TotalSale = total - header.Discount
	header.ProfitEstimate = profit
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateSaleHeader(ctx, header); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalSale, string(header.Payment), header.UpdatedAt); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.SaleDate, header.UserID, header.BranchID, header.TotalSale-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return sale.Item{}, apperror.New(http.StatusInternalServerError, "Update sale item failed", err.Error())
	}
	return s.Repo.FindSaleItemByID(ctx, id)
}

func (s Service) DeleteItem(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindSaleItemByID(ctx, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete sale item failed", "item not found")
	}
	header, err := s.Repo.FindSaleByID(ctx, branchID, item.SaleID)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete sale item failed", "sale not found")
	}
	oldTotal := header.TotalSale
	oldProfit := header.ProfitEstimate
	prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	prod.Stock += item.Qty
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	if err := s.Repo.DeleteSaleItem(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	items, err := s.Repo.FindSaleItems(ctx, item.SaleID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	var total int
	var profit int
	for _, line := range items {
		total += line.SubTotal
		profit += line.SubTotal - (line.Price * line.Qty)
	}
	header.TotalSale = total - header.Discount
	header.ProfitEstimate = profit
	header.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateSaleHeader(ctx, header); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalSale, string(header.Payment), header.UpdatedAt); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.SaleDate, header.UserID, header.BranchID, header.TotalSale-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale item failed", err.Error())
	}
	return nil
}

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	items, err := s.Repo.ListSales(ctx, branchID, sale.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sales excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Sales")
	sheet := "Sales"
	f.SetCellValue(sheet, "A1", "DATA PENJUALAN")
	headers := []string{"ID", "MEMBER", "DATE", "TOTAL", "DISCOUNT", "PAYMENT", "CASHIER"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.MemberName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.SaleDate.Format("02/01/2006"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.TotalSale)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Discount)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Payment)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), item.Cashier)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sales excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("sales-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	items, err := s.Repo.ListSales(ctx, branchID, sale.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sales pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("PENJUALAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "PENJUALAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	widths := []float64{30, 55, 35, 35, 25, 40, 57}
	headers := []string{"ID", "MEMBER", "DATE", "TOTAL", "DISC", "PAYMENT", "CASHIER"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items.Items {
		values := []string{item.ID, item.MemberName, item.SaleDate.Format("02/01/2006"), fmt.Sprintf("%d", item.TotalSale), fmt.Sprintf("%d", item.Discount), string(item.Payment), item.Cashier}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sales pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("PENJUALAN-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportItemsExcel(ctx context.Context, branchID, saleID string) ([]byte, string, error) {
	items, err := s.ListItems(ctx, branchID, saleID)
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Sale Items")
	sheet := "Sale Items"
	f.SetCellValue(sheet, "A1", "DETAIL PENJUALAN")
	headers := []string{"ID", "PRODUCT", "UNIT", "PRICE", "QTY", "SUB TOTAL"}
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
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale items excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-PENJUALAN-%s-%s.xlsx", saleID, time.Now().Format("20060102150405")), nil
}

func (s Service) ExportItemsPDF(ctx context.Context, branchID, saleID string) ([]byte, string, error) {
	items, err := s.ListItems(ctx, branchID, saleID)
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("DETAIL PENJUALAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "DETAIL PENJUALAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	widths := []float64{30, 95, 40, 35, 25, 52}
	headers := []string{"ID", "PRODUCT", "UNIT", "PRICE", "QTY", "SUB TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items {
		values := []string{item.ID, item.ProductName, item.UnitName, fmt.Sprintf("%d", item.Price), fmt.Sprintf("%d", item.Qty), fmt.Sprintf("%d", item.SubTotal)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export sale items pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("DETAIL-PENJUALAN-%s-%s.pdf", saleID, time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) CreateTransaction(ctx context.Context, branchID, userID, defaultMember, subscriptionType string, req sale.CreateSaleRequest) (sale.Sale, []sale.Item, error) {
	now := s.Clock.Now()
	payment := common.PaymentStatus(req.Sale.Payment)
	if payment == "" {
		payment = common.PaymentCash
	}
	memberID := req.Sale.MemberID
	if memberID == "" {
		memberID = defaultMember
	}
	saleEntity := sale.Sale{
		ID:        s.IDs.New("SAL"),
		MemberID:  memberID,
		UserID:    userID,
		BranchID:  branchID,
		Payment:   payment,
		Discount:  req.Sale.Discount,
		SaleDate:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	items := make([]sale.Item, 0, len(req.SaleItems))
	if err := s.Repo.WithinTransactionSale(ctx, func(repo ports.SaleTxRepository) error {
		var totalSale int
		var profitEstimate int
		for _, input := range req.SaleItems {
			prod, err := repo.FindProduct(ctx, input.ProductID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Product with ID %s not found", input.ProductID), err)
			}
			if prod.Stock < input.Qty {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", prod.Name, prod.Stock, input.Qty), nil)
			}
			serverPrice := prod.SalesPrice
			subtotal := serverPrice * input.Qty
			items = append(items, sale.Item{
				ID:        s.IDs.New("SIT"),
				SaleID:    saleEntity.ID,
				ProductID: input.ProductID,
				Price:     serverPrice,
				Qty:       input.Qty,
				SubTotal:  subtotal,
			})
			prod.Stock -= input.Qty
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update stock for product %s", prod.Name), err)
			}
			totalSale += subtotal
			profitEstimate += (serverPrice - prod.PurchasePrice) * input.Qty
		}
		saleEntity.TotalSale = totalSale - saleEntity.Discount
		saleEntity.ProfitEstimate = profitEstimate
		if err := repo.CreateSale(ctx, saleEntity); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create sale", err)
		}
		if err := repo.CreateSaleItems(ctx, items); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create sale items", err)
		}
		if err := repo.CreateTransactionReport(ctx, saleEntity.ID, "sale", userID, branchID, saleEntity.TotalSale, string(saleEntity.Payment), now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create transaction report", err)
		}
		if err := repo.UpsertDailyProfit(ctx, saleEntity.SaleDate, userID, branchID, saleEntity.TotalSale, saleEntity.ProfitEstimate, now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to update daily profit report", err)
		}
		if subscriptionType == "quota" {
			br, err := repo.FindBranch(ctx, branchID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Branch with ID %s not found", branchID), err)
			}
			if br.Quota <= 0 {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("No quota available for branch %s", br.BranchName), nil)
			}
			if err := repo.UpdateBranchQuota(ctx, branchID, br.Quota-1); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update quota for branch %s", br.BranchName), err)
			}
		}
		if memberID != "" && memberID != defaultMember {
			m, err := repo.FindMember(ctx, memberID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Member with ID %s not found", memberID), err)
			}
			mc, err := repo.FindMemberCategory(ctx, m.MemberCategoryID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Member category with ID %s not found for member %s", m.MemberCategoryID, m.ID), err)
			}
			if mc.PointsConversionRate > 0 {
				points := saleEntity.TotalSale / mc.PointsConversionRate
				if err := repo.UpdateMemberPoints(ctx, m.ID, m.Points+points); err != nil {
					return apperror.New(http.StatusInternalServerError, "Failed to update member points", err)
				}
			}
		}
		return nil
	}); err != nil {
		return sale.Sale{}, nil, err
	}
	return saleEntity, items, nil
}
