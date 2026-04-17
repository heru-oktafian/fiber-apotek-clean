package console

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	blue   = "\033[34m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	white  = "\033[97m"
)

func StartupBanner(appName, version, host, port string, handlers int, prefork bool) string {
	if appName == "" {
		appName = "Rest API Apotek"
	}
	if version == "" {
		version = "dev"
	}
	if host == "" {
		host = "0.0.0.0"
	}
	processes := 1
	preforkText := "Disabled"
	if prefork {
		preforkText = "Enabled"
	}

	header := fmt.Sprintf("%s%s%s | %s%s%s", bold, white, appName, yellow, version, reset)
	listen := fmt.Sprintf("http://%s:%s", host, port)
	statsRowOne := []boxStat{
		{Label: "Handlers", Value: fmt.Sprintf("%d", handlers), LabelColor: yellow},
		{Label: "Processes", Value: fmt.Sprintf("%d", processes), LabelColor: blue},
	}
	statsRowTwo := []boxStat{
		{Label: "PID", Value: fmt.Sprintf("%d", os.Getpid()), LabelColor: red},
		{Label: "Prefork", Value: preforkText, LabelColor: cyan},
	}

	lines := []string{
		boxLine("╔", "╗", 78),
		boxCentered(header, 78),
		boxCentered(fmt.Sprintf("%s%s%s", gray, listen, reset), 78),
		boxCentered(fmt.Sprintf("%s(bound on host %s and port %s)%s", gray, host, port, reset), 78),
		boxLine("╠", "╣", 78),
		boxStatRow(statsRowOne, 78),
		boxStatRow(statsRowTwo, 78),
		boxLine("╚", "╝", 78),
	}
	return strings.Join(lines, "\n")
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		started := time.Now()
		err := c.Next()
		latency := time.Since(started)
		status := c.Response().StatusCode()
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s %s[%3d]%s %s%-6s%s %s%-24s%s %s%-15s%s %s(%s)%s\n",
			timestamp,
			colorForStatus(status), status, reset,
			colorForMethod(c.Method()), c.Method(), reset,
			white, c.Path(), reset,
			gray, c.IP(), reset,
			gray, latency.Truncate(time.Microsecond), reset,
		)
		return err
	}
}

func boxLine(left, right string, width int) string {
	return left + strings.Repeat("═", width-2) + right
}

func boxCentered(text string, width int) string {
	innerWidth := width - 4
	padding := innerWidth - printableLen(text)
	if padding < 0 {
		padding = 0
	}
	leftPad := padding / 2
	rightPad := padding - leftPad
	return fmt.Sprintf("║ %s%s%s ║", strings.Repeat(" ", leftPad), text, strings.Repeat(" ", rightPad))
}

type boxStat struct {
	Label      string
	Value      string
	LabelColor string
}

func boxStatRow(stats []boxStat, width int) string {
	innerWidth := width - 4
	columnGap := 4
	columnWidth := (innerWidth - columnGap) / 2
	parts := make([]string, 0, len(stats))
	for _, stat := range stats {
		parts = append(parts, formatBoxStat(stat, columnWidth))
	}
	joined := strings.Join(parts, strings.Repeat(" ", columnGap))
	padding := innerWidth - printableLen(joined)
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("║ %s%s ║", joined, strings.Repeat(" ", padding))
}

func formatBoxStat(stat boxStat, width int) string {
	label := fmt.Sprintf("%s%-10s%s", stat.LabelColor, stat.Label, reset)
	value := fmt.Sprintf("%s%s%s", white, stat.Value, reset)
	plain := stat.Label + stat.Value
	spacing := width - printableLen(label) - printableLen(value)
	if spacing < 1 {
		spacing = 1
	}
	_ = plain
	return label + strings.Repeat(" ", spacing) + value
}

func printableLen(value string) int {
	length := 0
	inEscape := false
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch == 0x1b {
			inEscape = true
			continue
		}
		if inEscape {
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

func colorForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return blue
	case status >= 400 && status < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case fiber.MethodGet:
		return cyan
	case fiber.MethodPost:
		return green
	case fiber.MethodPut, fiber.MethodPatch:
		return yellow
	case fiber.MethodDelete:
		return red
	default:
		return white
	}
}
