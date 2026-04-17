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
