package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/salereturn"
	sharedresponse "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	salereturnusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/salereturn"
)

func currentSaleReturnUserID(claims auth.Claims) string {
	return claims.Subject
}

type SaleReturnHandler struct {
	Service salereturnusecase.Service
}

func (h SaleReturnHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	result, err := h.Service.List(c.Context(), claims.BranchID, salereturn.ListRequest{
		Search: c.Query("search"),
		Month:  c.Query("month"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
	})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSONWithMeta(c, fiber.StatusOK, "Data retur penjualan berhasil diambil", result.Items, result.Meta)
}

func (h SaleReturnHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	result, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Retur penjualan berhasil diambil", result)
}

func (h SaleReturnHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req salereturn.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}
	result, err := h.Service.Create(c.Context(), claims.BranchID, currentSaleReturnUserID(claims), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Transaksi retur penjualan berhasil dibuat", result)
}
