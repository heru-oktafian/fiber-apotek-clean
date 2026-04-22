package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/suppliercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	suppliercategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/suppliercategory"
)

type SupplierCategoryHandler struct {
	Service suppliercategoryusecase.Service
}

func (h SupplierCategoryHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, suppliercategory.ListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Supplier Category retrieved successfully", result.Items, result.Meta)
}

func (h SupplierCategoryHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, uint(id))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Get supplier category success", item)
}

func (h SupplierCategoryHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req suppliercategory.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "Create supplier category success", item)
}

func (h SupplierCategoryHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	var req suppliercategory.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, uint(id), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Update supplier category success", item)
}

func (h SupplierCategoryHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	if err := h.Service.Delete(c.Context(), claims.BranchID, uint(id)); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Delete supplier category success", fiber.Map{"id": id})
}

func (h SupplierCategoryHandler) Combo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.Combo(c.Context(), claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data berhasil ditemukan", items)
}
