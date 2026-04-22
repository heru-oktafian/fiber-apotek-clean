package unit

import (
	"context"
	"net/http"
	"strings"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type MasterService struct {
	Units ports.UnitRepository
	IDs   ports.IDGenerator
}

func (s MasterService) List(ctx context.Context, branchID string, req domain.MasterUnitListRequest) (domain.MasterUnitListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Units.ListMasterUnits(ctx, branchID, req)
	if err != nil {
		return domain.MasterUnitListResult{}, apperror.New(http.StatusInternalServerError, "Get units failed", err.Error())
	}
	return items, nil
}

func (s MasterService) GetByID(ctx context.Context, branchID, id string) (domain.MasterUnit, error) {
	item, err := s.Units.FindMasterUnitByID(ctx, id, branchID)
	if err != nil {
		return domain.MasterUnit{}, apperror.New(http.StatusNotFound, "Get unit failed", err.Error())
	}
	return item, nil
}

func (s MasterService) Create(ctx context.Context, branchID string, req domain.MasterUnitCreateRequest) (domain.MasterUnit, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return domain.MasterUnit{}, apperror.New(http.StatusBadRequest, "Create unit failed", "name is required")
	}
	item := domain.MasterUnit{ID: s.IDs.New("UNT"), Name: name, BranchID: branchID}
	if err := s.Units.CreateMasterUnit(ctx, item); err != nil {
		return domain.MasterUnit{}, apperror.New(http.StatusInternalServerError, "Create unit failed", err.Error())
	}
	return item, nil
}

func (s MasterService) Update(ctx context.Context, branchID, id string, req domain.MasterUnitCreateRequest) (domain.MasterUnit, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return domain.MasterUnit{}, apperror.New(http.StatusBadRequest, "Update unit failed", "name is required")
	}
	if _, err := s.Units.FindMasterUnitByID(ctx, id, branchID); err != nil {
		return domain.MasterUnit{}, apperror.New(http.StatusNotFound, "Update unit failed", "unit not found")
	}
	item := domain.MasterUnit{ID: id, Name: name, BranchID: branchID}
	if err := s.Units.UpdateMasterUnit(ctx, item); err != nil {
		return domain.MasterUnit{}, apperror.New(http.StatusInternalServerError, "Update unit failed", err.Error())
	}
	return s.Units.FindMasterUnitByID(ctx, id, branchID)
}

func (s MasterService) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Units.FindMasterUnitByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete unit failed", "unit not found")
	}
	if err := s.Units.DeleteMasterUnit(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete unit failed", err.Error())
	}
	return nil
}

func (s MasterService) Combo(ctx context.Context, branchID, search string) ([]domain.MasterUnitComboItem, error) {
	items, err := s.Units.GetMasterUnitCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get unit combo failed", err.Error())
	}
	return items, nil
}
