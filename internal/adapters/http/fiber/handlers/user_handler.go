package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	userusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/user"
)

type UserHandler struct {
	Service userusecase.Service
}

func (h UserHandler) Create(c *fiber.Ctx) error {
	var req user.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	result, err := h.Service.Create(c.Context(), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusCreated, "User berhasil dibuat", result)
}

func (h UserHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), user.ListRequest{
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Data berhasil diambil", result.Items, result.Meta)
}

func (h UserHandler) Detail(c *fiber.Ctx) error {
	result, err := h.Service.Detail(c.Context(), c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Data berhasil ditemukan", result)
}

func (h UserHandler) Update(c *fiber.Ctx) error {
	var req user.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	result, err := h.Service.Update(c.Context(), c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "User berhasil diupdate", result)
}
