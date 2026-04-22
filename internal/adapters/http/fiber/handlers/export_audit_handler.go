package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
	firststockusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/firststock"
)

type ExportAuditHandler struct {
	FirstStocks firststockusecase.Service
}

func (h ExportAuditHandler) FirstStocksExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.FirstStocks.ExportExcel(c.Context(), branchID, f.Month)
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportAuditHandler) FirstStocksPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.FirstStocks.ExportPDF(c.Context(), branchID, f.Month)
	}, "application/pdf")
}

func (h ExportAuditHandler) FirstStockItemsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.FirstStocks.ExportItemsExcel(c.Context(), branchID, c.Query("first_stock_id"))
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportAuditHandler) FirstStockItemsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.FirstStocks.ExportItemsPDF(c.Context(), branchID, c.Query("first_stock_id"))
	}, "application/pdf")
}

func (h ExportAuditHandler) send(c *fiber.Ctx, fn func(branchID string, filters exportshared.Filters) ([]byte, string, error), contentType string) error {
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
