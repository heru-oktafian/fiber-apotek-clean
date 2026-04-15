package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	productusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/product"
)

type ProductHandler struct {
	Service productusecase.Service
}

func (h ProductHandler) Create(c *fiber.Ctx) error {
	var req product.Product
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	claims := c.Locals("claims").(auth.Claims)
	created, err := h.Service.Create(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Resource created successfully", created)
}

func (h ProductHandler) SaleCombo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.SaleCombo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Combo Products retrieved successfully", items)
}

func (h ProductHandler) PurchaseCombo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.PurchaseCombo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Combo Purchase Products retrieved successfully", items)
}

func (h ProductHandler) OpnameCombo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.OpnameCombo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data Combobox ditemukan", items)
}
