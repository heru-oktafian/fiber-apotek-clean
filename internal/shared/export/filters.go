package export

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Filters struct {
	Search string
	Month  string
	ID     string
}

func ParseFilters(c *fiber.Ctx) Filters {
	month := strings.TrimSpace(c.Query("month"))
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return Filters{
		Search: strings.TrimSpace(c.Query("search")),
		Month:  month,
		ID:     strings.TrimSpace(c.Params("id")),
	}
}
