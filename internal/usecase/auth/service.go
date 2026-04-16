package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Users     ports.UserRepository
	Branches  ports.BranchRepository
	Passwords ports.PasswordComparer
	Tokens    ports.TokenManager
	Blacklist ports.TokenBlacklist
	Clock     ports.Clock
}

func (s Service) Login(ctx context.Context, req auth.LoginRequest) (auth.LoginResult, error) {
	usr, err := s.Users.FindActiveByUsername(ctx, req.Username)
	if err != nil {
		return auth.LoginResult{}, apperror.New(http.StatusUnauthorized, "Login failed", "User is not active, call admin to activated your account !")
	}
	if err := s.Passwords.Compare(usr.Password, req.Password); err != nil {
		return auth.LoginResult{}, apperror.New(http.StatusUnauthorized, "Login failed", "Invalid username or password")
	}
	expiresAt := s.Clock.Now().Add(5 * time.Minute)
	token, err := s.Tokens.GenerateLoginToken(usr, expiresAt)
	if err != nil {
		return auth.LoginResult{}, apperror.New(http.StatusInternalServerError, "Login failed", "Failed to generate token")
	}
	return auth.LoginResult{Token: token}, nil
}

func (s Service) ListBranches(ctx context.Context, bearerToken string) ([]auth.UserBranch, error) {
	claims, _, err := s.parseBearerToken(bearerToken, "Insert valid token to access this endpoint!")
	if err != nil {
		return nil, err
	}
	items, err := s.Branches.ListUserBranches(ctx, claims.Subject)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}
	return items, nil
}

func (s Service) SetBranch(ctx context.Context, bearerToken string, req auth.BranchSelectionRequest) (string, error) {
	raw := strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	claims, exp, err := s.parseBearerToken(bearerToken, "Insert valid token to access this endpoint!")
	if err != nil {
		return "", err
	}
	hasBranch, err := s.Branches.UserHasBranch(ctx, claims.Subject, req.BranchID)
	if err != nil || !hasBranch {
		return "", apperror.New(http.StatusForbidden, "Invalid branch ID", "Branch not associated with this user!")
	}
	usr, err := s.Users.FindByID(ctx, claims.Subject)
	if err != nil {
		return "", apperror.New(http.StatusInternalServerError, "Failed to set branch", "Unable to retrieve user role")
	}
	br, err := s.Branches.FindBranchByID(ctx, req.BranchID)
	if err != nil {
		return "", apperror.New(http.StatusInternalServerError, "Failed to set branch", "Unable to retrieve branch details")
	}
	newToken, err := s.Tokens.GenerateBranchToken(auth.Claims{
		Subject:          usr.ID,
		Name:             usr.Name,
		BranchID:         br.ID,
		UserRole:         string(usr.Role),
		DefaultMember:    br.DefaultMemberID,
		Quota:            br.Quota,
		SubscriptionType: br.SubscriptionType,
		RealAsset:        br.RealAsset,
	}, s.Clock.Now().Add(8*time.Hour))
	if err != nil {
		return "", apperror.New(http.StatusInternalServerError, "Failed to set branch", "Failed to generate new token")
	}
	if err := s.Blacklist.Blacklist(ctx, raw, time.Until(exp)); err != nil {
		return "", apperror.New(http.StatusInternalServerError, "Failed to set branch", "Failed to blacklist old token")
	}
	return newToken, nil
}

func (s Service) Profile(ctx context.Context, bearerToken string) (auth.Profile, error) {
	claims, _, err := s.parseBearerToken(bearerToken, "Insert valid token to access this endpoint!")
	if err != nil {
		return auth.Profile{}, err
	}
	item, err := s.Branches.FindProfile(ctx, claims.Subject, claims.BranchID)
	if err != nil {
		return auth.Profile{}, apperror.New(http.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}
	if !itemHasProfile(item) {
		return auth.Profile{}, apperror.New(http.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}
	return item, nil
}

func (s Service) Logout(ctx context.Context, bearerToken string) error {
	raw := strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	if raw == "" {
		return apperror.New(http.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint !")
	}
	_, exp, err := s.Tokens.Parse(raw)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Logout failed", "Failed to blacklist token")
	}
	if err := s.Blacklist.Blacklist(ctx, raw, time.Until(exp)); err != nil {
		return apperror.New(http.StatusInternalServerError, "Logout failed", "Failed to blacklist token")
	}
	return nil
}

func (s Service) parseBearerToken(bearerToken string, missingMessage string) (auth.Claims, time.Time, error) {
	raw := strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	if raw == "" {
		return auth.Claims{}, time.Time{}, apperror.New(http.StatusUnauthorized, "Missing token", missingMessage)
	}
	claims, exp, err := s.Tokens.Parse(raw)
	if err != nil {
		return auth.Claims{}, time.Time{}, apperror.New(http.StatusUnauthorized, "Invalid token", "Try to login again!")
	}
	return claims, exp, nil
}

func itemHasProfile(item auth.Profile) bool {
	return item.UserID != "" && item.BranchID != ""
}
