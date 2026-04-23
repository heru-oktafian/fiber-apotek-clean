package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/buyreturn"
	sharedresponse "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	buyreturnusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/buyreturn"
)

func currentUserID(claims auth.Claims) string {
	return claims.Subject
}

type BuyReturnHandler struct {
	Service buyreturnusecase.Service
}

func (h BuyReturnHandler) PurchaseSources(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.ListPurchaseSources(c.Context(), claims.BranchID, c.Query("search"), c.Query("month"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Data pembelian berhasil diambil", items)
}

func (h BuyReturnHandler) ReturnableItems(c *fiber.Ctx) error {
	items, err := h.Service.ListReturnableItems(c.Context(), c.Query("purchase_id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Data item retur ditemukan", items)
}

func (h BuyReturnHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	result, err := h.Service.List(c.Context(), claims.BranchID, buyreturn.ListRequest{
		Search: c.Query("search"),
		Month:  c.Query("month"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
	})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSONWithMeta(c, fiber.StatusOK, "Data retur pembelian berhasil diambil", result.Items, result.Meta)
}

func (h BuyReturnHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	result, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Retur pembelian berhasil diambil", result)
}

func (h BuyReturnHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req buyreturn.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}
	result, err := h.Service.Create(c.Context(), claims.BranchID, currentUserID(claims), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return sharedresponse.JSON(c, fiber.StatusOK, "Transaksi retur pembelian berhasil dibuat", result)
}
