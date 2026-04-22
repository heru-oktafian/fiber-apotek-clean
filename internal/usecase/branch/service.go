package branch

import (
	"context"
	"net/http"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Branches ports.BranchRepository
	IDs      ports.IDGenerator
}

func (s Service) List(ctx context.Context, req branch.ListRequest) (branch.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Branches.ListBranches(ctx, req)
	if err != nil {
		return branch.ListResult{}, apperror.New(http.StatusInternalServerError, "Get all branch failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, id string) (branch.Branch, error) {
	item, err := s.Branches.FindBranchByID(ctx, id)
	if err != nil {
		return branch.Branch{}, apperror.New(http.StatusNotFound, "Get branch failed", err.Error())
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, req branch.CreateRequest) (branch.Branch, error) {
	req.BranchName = strings.TrimSpace(req.BranchName)
	req.BranchStatus = strings.TrimSpace(req.BranchStatus)
	if req.BranchName == "" {
		return branch.Branch{}, apperror.New(http.StatusBadRequest, "Create branch failed", "branch_name is required")
	}
	if req.BranchStatus == "" {
		req.BranchStatus = "active"
	}

	item := branch.Branch{
		ID:               s.IDs.New("BRC"),
		BranchName:       req.BranchName,
		Address:          strings.TrimSpace(req.Address),
		Phone:            strings.TrimSpace(req.Phone),
		Email:            strings.TrimSpace(req.Email),
		SIAID:            strings.TrimSpace(req.SIAID),
		SIAName:          strings.TrimSpace(req.SIAName),
		PSAID:            strings.TrimSpace(req.PSAID),
		PSAName:          strings.TrimSpace(req.PSAName),
		SIPA:             strings.TrimSpace(req.SIPA),
		SIPAName:         strings.TrimSpace(req.SIPAName),
		APINGID:          strings.TrimSpace(req.APINGID),
		APINGName:        strings.TrimSpace(req.APINGName),
		BankName:         strings.TrimSpace(req.BankName),
		AccountName:      strings.TrimSpace(req.AccountName),
		AccountNumber:    strings.TrimSpace(req.AccountNumber),
		TaxPercentage:    req.TaxPercentage,
		JournalMethod:    strings.TrimSpace(req.JournalMethod),
		BranchStatus:     req.BranchStatus,
		LicenseDate:      strings.TrimSpace(req.LicenseDate),
		DefaultMemberID:  strings.TrimSpace(req.DefaultMemberID),
		SubscriptionType: strings.TrimSpace(req.SubscriptionType),
		Quota:            req.Quota,
		RealAsset:        strings.TrimSpace(req.RealAsset),
	}

	if err := s.Branches.CreateBranch(ctx, item); err != nil {
		return branch.Branch{}, apperror.New(http.StatusInternalServerError, "Create branch failed", err.Error())
	}
	return item, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(http.StatusBadRequest, "Delete branch failed", "branch id is required")
	}
	if _, err := s.Branches.FindBranchByID(ctx, id); err != nil {
		return apperror.New(http.StatusNotFound, "Delete branch failed", "branch not found")
	}
	hasUsers, err := s.Branches.BranchHasUsers(ctx, id)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete branch failed", err.Error())
	}
	if hasUsers {
		return apperror.New(http.StatusBadRequest, "Delete branch failed", "branch masih terhubung ke user")
	}
	if err := s.Branches.DeleteBranch(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete branch failed", err.Error())
	}
	return nil
}
