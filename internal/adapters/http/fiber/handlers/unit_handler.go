package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	unitusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/unit"
)

type UnitHandler struct {
	Service unitusecase.MasterService
}

func (h UnitHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, domain.MasterUnitListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Units retrieved successfully", result.Items, result.Meta)
}

func (h UnitHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Get unit success", item)
}

func (h UnitHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req domain.MasterUnitCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "Create unit success", item)
}

func (h UnitHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req domain.MasterUnitCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Update unit success", item)
}

func (h UnitHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Delete unit success", fiber.Map{"id": c.Params("id")})
}

func (h UnitHandler) Combo(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.Combo(c.Context(), claims.BranchID, c.Query("search"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data berhasil ditemukan", items)
}
