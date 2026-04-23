package unit

import (
	"context"
	"net/http"
	"strings"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type ConversionService struct {
	Units ports.UnitRepository
	IDs   ports.IDGenerator
}

func (s ConversionService) List(ctx context.Context, branchID string, req domain.ConversionListRequest) (domain.ConversionListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Units.ListConversions(ctx, branchID, req)
	if err != nil {
		return domain.ConversionListResult{}, apperror.New(http.StatusInternalServerError, "Get unit conversions failed", err.Error())
	}
	return items, nil
}

func (s ConversionService) GetByID(ctx context.Context, branchID, id string) (domain.ConversionMaster, error) {
	item, err := s.Units.FindConversionByID(ctx, id, branchID)
	if err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Get unit conversion failed", "unit conversion not found")
	}
	return item, nil
}

func (s ConversionService) Create(ctx context.Context, branchID string, req domain.ConversionCreateRequest) (domain.ConversionMaster, error) {
	productID := strings.TrimSpace(req.ProductID)
	initID := strings.TrimSpace(req.InitID)
	finalID := strings.TrimSpace(req.FinalID)
	if productID == "" || initID == "" || finalID == "" {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "product_id, init_id, and final_id are required")
	}
	if req.ValueConv <= 0 {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "value_conv must be greater than zero")
	}
	if initID == finalID {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Create unit conversion failed", "init_id and final_id must be different")
	}
	if _, err := s.Units.FindProductByID(ctx, productID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "product not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, initID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "initial unit not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, finalID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Create unit conversion failed", "final unit not found")
	}
	if _, err := s.Units.FindConversion(ctx, productID, initID, finalID, branchID); err == nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusConflict, "Create unit conversion failed", "unit conversion already exists")
	}
	item := domain.ConversionMaster{ID: s.IDs.New("UNC"), ProductID: productID, InitID: initID, FinalID: finalID, ValueConv: req.ValueConv, BranchID: branchID}
	if err := s.Units.CreateConversion(ctx, item); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusInternalServerError, "Create unit conversion failed", err.Error())
	}
	return s.Units.FindConversionByID(ctx, item.ID, branchID)
}

func (s ConversionService) Update(ctx context.Context, branchID, id string, req domain.ConversionCreateRequest) (domain.ConversionMaster, error) {
	productID := strings.TrimSpace(req.ProductID)
	initID := strings.TrimSpace(req.InitID)
	finalID := strings.TrimSpace(req.FinalID)
	if productID == "" || initID == "" || finalID == "" {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "product_id, init_id, and final_id are required")
	}
	if req.ValueConv <= 0 {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "value_conv must be greater than zero")
	}
	if initID == finalID {
		return domain.ConversionMaster{}, apperror.New(http.StatusBadRequest, "Update unit conversion failed", "init_id and final_id must be different")
	}
	if _, err := s.Units.FindConversionByID(ctx, id, branchID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "unit conversion not found")
	}
	if _, err := s.Units.FindProductByID(ctx, productID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "product not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, initID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "initial unit not found")
	}
	if _, err := s.Units.FindUnitByID(ctx, finalID); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusNotFound, "Update unit conversion failed", "final unit not found")
	}
	item := domain.ConversionMaster{ID: id, ProductID: productID, InitID: initID, FinalID: finalID, ValueConv: req.ValueConv, BranchID: branchID}
	if err := s.Units.UpdateConversion(ctx, item); err != nil {
		return domain.ConversionMaster{}, apperror.New(http.StatusInternalServerError, "Update unit conversion failed", err.Error())
	}
	return s.Units.FindConversionByID(ctx, id, branchID)
}

func (s ConversionService) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Units.FindConversionByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete unit conversion failed", "unit conversion not found")
	}
	if err := s.Units.DeleteConversion(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete unit conversion failed", err.Error())
	}
	return nil
}
