package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
	buyreturnusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/buyreturn"
	salereturnusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/salereturn"
)

type ExportReturnHandler struct {
	BuyReturns  buyreturnusecase.Service
	SaleReturns salereturnusecase.Service
}

func (h ExportReturnHandler) BuyReturnsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.BuyReturns.ExportExcel(c.Context(), branchID, f.Month)
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportReturnHandler) BuyReturnsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.BuyReturns.ExportPDF(c.Context(), branchID, f.Month)
	}, "application/pdf")
}

func (h ExportReturnHandler) BuyReturnItemsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.BuyReturns.ExportItemsExcel(c.Context(), branchID, c.Query("buy_return_id"))
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportReturnHandler) BuyReturnItemsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.BuyReturns.ExportItemsPDF(c.Context(), branchID, c.Query("buy_return_id"))
	}, "application/pdf")
}

func (h ExportReturnHandler) SaleReturnsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.SaleReturns.ExportExcel(c.Context(), branchID, f.Month)
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportReturnHandler) SaleReturnsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.SaleReturns.ExportPDF(c.Context(), branchID, f.Month)
	}, "application/pdf")
}

func (h ExportReturnHandler) SaleReturnItemsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.SaleReturns.ExportItemsExcel(c.Context(), branchID, c.Query("sale_return_id"))
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportReturnHandler) SaleReturnItemsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.SaleReturns.ExportItemsPDF(c.Context(), branchID, c.Query("sale_return_id"))
	}, "application/pdf")
}

func (h ExportReturnHandler) send(c *fiber.Ctx, fn func(branchID string, filters exportshared.Filters) ([]byte, string, error), contentType string) error {
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
