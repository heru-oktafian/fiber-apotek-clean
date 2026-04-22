package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	branchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/branch"
)

type BranchHandler struct {
	Service branchusecase.Service
}

func (h BranchHandler) Create(c *fiber.Ctx) error {
	var req branch.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "Create branch success", item)
}

func (h BranchHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), branch.ListRequest{
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Get all branch success", result.Items, result.Meta)
}

func (h BranchHandler) GetByID(c *fiber.Ctx) error {
	item, err := h.Service.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Get branch success", item)
}

func (h BranchHandler) Delete(c *fiber.Ctx) error {
	if err := h.Service.Delete(c.Context(), c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Delete branch success", fiber.Map{"id": c.Params("id")})
}
