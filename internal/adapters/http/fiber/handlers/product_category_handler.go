package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/productcategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	productcategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/productcategory"
)

type ProductCategoryHandler struct {
	Service productcategoryusecase.Service
}

func (h ProductCategoryHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, productcategory.ListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Product Categories retrieved successfully", result.Items, result.Meta)
}

func (h ProductCategoryHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, uint(id))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Get product category success", item)
}

func (h ProductCategoryHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req productcategory.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "Create product category success", item)
}

func (h ProductCategoryHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	var req productcategory.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, uint(id), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Update product category success", item)
}

func (h ProductCategoryHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return presenter.Handle(c, fiber.ErrBadRequest)
	}
	if err := h.Service.Delete(c.Context(), claims.BranchID, uint(id)); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Delete product category success", fiber.Map{"id": id})
}

func (h ProductCategoryHandler) Combo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.Combo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data berhasil ditemukan", items)
}
