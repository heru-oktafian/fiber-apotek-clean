package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/supplier"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	supplierusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/supplier"
)

type SupplierHandler struct {
	Service supplierusecase.Service
}

func (h SupplierHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, supplier.ListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Suppliers retrieved successfully", result.Items, result.Meta)
}

func (h SupplierHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Get supplier success", item)
}

func (h SupplierHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req supplier.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "Create supplier success", item)
}

func (h SupplierHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req supplier.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Update supplier success", item)
}

func (h SupplierHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Delete supplier success", fiber.Map{"id": c.Params("id")})
}

func (h SupplierHandler) Combo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.Combo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data berhasil ditemukan", items)
}
