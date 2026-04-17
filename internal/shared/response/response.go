package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func JSON(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  statusText(status),
		"message": message,
		"data":    data,
	})
}

func JSONWithMeta(c *fiber.Ctx, status int, message string, data any, meta any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  statusText(status),
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *fiber.Ctx, status int, message string, detail any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  statusText(status),
		"message": message,
		"error":   detail,
	})
}

func statusText(status int) string {
	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		return "success"
	}
	return "error"
}
