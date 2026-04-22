package export

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FileFormat string

const (
	FormatExcel FileFormat = "xlsx"
	FormatPDF   FileFormat = "pdf"
)

func SafeFilename(base string, format FileFormat) string {
	clean := strings.TrimSpace(strings.ToLower(base))
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "/", "-")
	clean = strings.ReplaceAll(clean, "_", "-")
	if clean == "" {
		clean = "export"
	}
	return fmt.Sprintf("%s-%s.%s", clean, time.Now().Format("20060102-150405"), string(format))
}

func ContentType(format FileFormat) string {
	switch format {
	case FormatExcel:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case FormatPDF:
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func SendFile(c *fiber.Ctx, filename string, format FileFormat, data []byte) error {
	c.Set("Content-Type", ContentType(format))
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	return c.Status(fiber.StatusOK).Send(data)
}
