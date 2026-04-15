package response

import "github.com/gofiber/fiber/v2"

func JSON(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

func Error(c *fiber.Ctx, status int, message string, detail any) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"error":   detail,
	})
}
