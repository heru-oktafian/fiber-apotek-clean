package presenter

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
)

func Handle(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*apperror.Error); ok {
		return response.Error(c, appErr.Code, appErr.Message, appErr.Detail)
	}
	return response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
}
