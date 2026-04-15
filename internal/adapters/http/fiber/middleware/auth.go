package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type TokenService interface {
	Parse(token string) (auth.Claims, time.Time, error)
}

type BlacklistService interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

func RequireAuth(tokens TokenService, blacklist BlacklistService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
		if raw == "" {
			return presenter.Handle(c, apperror.New(fiber.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint !"))
		}
		claims, _, err := tokens.Parse(raw)
		if err != nil {
			return presenter.Handle(c, apperror.New(fiber.StatusUnauthorized, "Invalid token", "Try to login again!"))
		}
		blocked, err := blacklist.IsBlacklisted(c.Context(), raw)
		if err == nil && blocked {
			return presenter.Handle(c, apperror.New(fiber.StatusUnauthorized, "Invalid token", "Try to login again!"))
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
