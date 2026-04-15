package jwt

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
)

type Service struct {
	Secret []byte
}

func (s Service) GenerateLoginToken(user user.User, expiresAt time.Time) (string, error) {
	claims := jwtv5.MapClaims{
		"sub": user.ID,
		"exp": expiresAt.Unix(),
	}
	return jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims).SignedString(s.Secret)
}

func (s Service) GenerateBranchToken(claims auth.Claims, expiresAt time.Time) (string, error) {
	mapClaims := jwtv5.MapClaims{
		"sub":               claims.Subject,
		"name":              claims.Name,
		"branch_id":         claims.BranchID,
		"user_role":         claims.UserRole,
		"default_member":    claims.DefaultMember,
		"quota":             claims.Quota,
		"subscription_type": claims.SubscriptionType,
		"real_asset":        claims.RealAsset,
		"exp":               expiresAt.Unix(),
	}
	return jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, mapClaims).SignedString(s.Secret)
}

func (s Service) Parse(token string) (auth.Claims, time.Time, error) {
	parsed, err := jwtv5.Parse(token, func(t *jwtv5.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil || !parsed.Valid {
		return auth.Claims{}, time.Time{}, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := parsed.Claims.(jwtv5.MapClaims)
	if !ok {
		return auth.Claims{}, time.Time{}, fmt.Errorf("invalid claims")
	}
	expUnix, _ := claims["exp"].(float64)
	result := auth.Claims{}
	if v, ok := claims["sub"].(string); ok {
		result.Subject = v
	}
	if v, ok := claims["name"].(string); ok {
		result.Name = v
	}
	if v, ok := claims["branch_id"].(string); ok {
		result.BranchID = v
	}
	if v, ok := claims["user_role"].(string); ok {
		result.UserRole = v
	}
	if v, ok := claims["default_member"].(string); ok {
		result.DefaultMember = v
	}
	if v, ok := claims["subscription_type"].(string); ok {
		result.SubscriptionType = v
	}
	if v, ok := claims["real_asset"].(string); ok {
		result.RealAsset = v
	}
	if v, ok := claims["quota"].(float64); ok {
		result.Quota = int(v)
	}
	return result, time.Unix(int64(expUnix), 0), nil
}
