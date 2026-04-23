package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
	duplicatereceiptusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/duplicatereceipt"
	purchaseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/purchase"
	saleusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/sale"
)

type ExportTransactionHandler struct {
	Purchases         purchaseusecase.Service
	Sales             saleusecase.Service
	DuplicateReceipts duplicatereceiptusecase.Service
}

func (h ExportTransactionHandler) PurchasesExcel(c *fiber.Ctx) error             { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Purchases.ExportExcel(c.Context(), branchID, f.Month) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) PurchasesPDF(c *fiber.Ctx) error               { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Purchases.ExportPDF(c.Context(), branchID, f.Month) }, "application/pdf") }
func (h ExportTransactionHandler) PurchaseItemsExcel(c *fiber.Ctx) error         { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Purchases.ExportItemsExcel(c.Context(), branchID, c.Query("purchase_id")) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) PurchaseItemsPDF(c *fiber.Ctx) error           { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Purchases.ExportItemsPDF(c.Context(), branchID, c.Query("purchase_id")) }, "application/pdf") }
func (h ExportTransactionHandler) SalesExcel(c *fiber.Ctx) error                 { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Sales.ExportExcel(c.Context(), branchID, f.Month) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) SalesPDF(c *fiber.Ctx) error                   { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Sales.ExportPDF(c.Context(), branchID, f.Month) }, "application/pdf") }
func (h ExportTransactionHandler) SaleItemsExcel(c *fiber.Ctx) error             { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Sales.ExportItemsExcel(c.Context(), branchID, c.Query("sale_id")) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) SaleItemsPDF(c *fiber.Ctx) error               { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.Sales.ExportItemsPDF(c.Context(), branchID, c.Query("sale_id")) }, "application/pdf") }
func (h ExportTransactionHandler) DuplicateReceiptsExcel(c *fiber.Ctx) error     { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.DuplicateReceipts.ExportExcel(c.Context(), branchID, f.Month) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) DuplicateReceiptsPDF(c *fiber.Ctx) error       { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.DuplicateReceipts.ExportPDF(c.Context(), branchID, f.Month) }, "application/pdf") }
func (h ExportTransactionHandler) DuplicateReceiptItemsExcel(c *fiber.Ctx) error { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.DuplicateReceipts.ExportItemsExcel(c.Context(), branchID, c.Query("duplicate_receipt_id")) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportTransactionHandler) DuplicateReceiptItemsPDF(c *fiber.Ctx) error   { return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) { return h.DuplicateReceipts.ExportItemsPDF(c.Context(), branchID, c.Query("duplicate_receipt_id")) }, "application/pdf") }

func (h ExportTransactionHandler) send(c *fiber.Ctx, fn func(branchID string, filters exportshared.Filters) ([]byte, string, error), contentType string) error {
	claims := c.Locals("claims").(auth.Claims)
	filters := exportshared.ParseFilters(c)
	data, filename, err := fn(claims.BranchID, filters)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}
