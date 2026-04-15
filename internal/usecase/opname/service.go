package opname

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/opname"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	"gorm.io/gorm"
)

type Service struct {
	Repo  ports.OpnameRepository
	Clock ports.Clock
	IDs   ports.IDGenerator
}

func (s Service) CreateHeader(ctx context.Context, branchID, userID string, req opname.CreateOpnameRequest) (opname.Opname, error) {
	parsedDate, err := time.Parse("2006-01-02", req.OpnameDate)
	if err != nil {
		return opname.Opname{}, apperror.New(http.StatusBadRequest, "Format tanggal tidak valid. Gunakan YYYY-MM-DD", err)
	}
	now := s.Clock.Now()
	entity := opname.Opname{
		ID:          s.IDs.New("OPN"),
		Description: req.Description,
		BranchID:    branchID,
		UserID:      userID,
		OpnameDate:  parsedDate,
		TotalOpname: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.Repo.CreateOpname(ctx, entity); err != nil {
		return opname.Opname{}, apperror.New(http.StatusInternalServerError, "Gagal membuat data opname", err)
	}
	return entity, nil
}

func (s Service) CreateItem(ctx context.Context, req opname.CreateOpnameItemRequest) (opname.Item, error) {
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return opname.Item{}, apperror.New(http.StatusBadRequest, "Format expired_date tidak valid. Gunakan YYYY-MM-DD", err)
	}

	var result opname.Item
	err = s.Repo.WithinOpnameTransaction(ctx, func(repo ports.OpnameTxRepository) error {
		header, err := repo.FindOpnameByID(ctx, req.OpnameID)
		if err != nil {
			return apperror.New(http.StatusNotFound, "Data opname tidak ditemukan", err)
		}
		prod, err := repo.FindProductByID(ctx, req.ProductID)
		if err != nil {
			return apperror.New(http.StatusNotFound, "Produk tidak ditemukan", err)
		}

		qtyExist := prod.Stock
		subTotalExist := qtyExist * prod.PurchasePrice
		subTotal := req.Qty * req.Price

		item := opname.Item{
			ID:            s.IDs.New("OPI"),
			OpnameID:      header.ID,
			ProductID:     prod.ID,
			ProductName:   prod.Name,
			Qty:           req.Qty,
			QtyExist:      qtyExist,
			Price:         req.Price,
			SubTotal:      subTotal,
			SubTotalExist: subTotalExist,
			ExpiredDate:   expiredDate,
		}

		existing, err := repo.FindOpnameItemByOpnameAndProduct(ctx, header.ID, prod.ID)
		if err == nil {
			existing.Qty = item.Qty
			existing.QtyExist = item.QtyExist
			existing.Price = item.Price
			existing.SubTotal = item.SubTotal
			existing.SubTotalExist = item.SubTotalExist
			existing.ExpiredDate = item.ExpiredDate
			existing.ProductName = item.ProductName
			if err := repo.UpdateOpnameItem(ctx, existing); err != nil {
				return apperror.New(http.StatusInternalServerError, "Gagal memperbarui item opname", err)
			}
			prod.Stock = req.Qty
			prod.PurchasePrice = req.Price
			prod.ExpiredDate = expiredDate
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, "Gagal memperbarui stok produk", err)
			}
			if _, err := repo.RecalculateOpnameTotal(ctx, header.ID); err != nil {
				return apperror.New(http.StatusInternalServerError, "Gagal menghitung ulang total opname", err)
			}
			result = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(http.StatusInternalServerError, "Gagal memeriksa item opname", err)
		}

		if err := repo.CreateOpnameItem(ctx, item); err != nil {
			return apperror.New(http.StatusInternalServerError, "Gagal membuat item opname", err)
		}
		prod.Stock = req.Qty
		prod.PurchasePrice = req.Price
		prod.ExpiredDate = expiredDate
		if err := repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Gagal memperbarui stok produk", err)
		}
		if _, err := repo.RecalculateOpnameTotal(ctx, header.ID); err != nil {
			return apperror.New(http.StatusInternalServerError, "Gagal menghitung ulang total opname", err)
		}
		result = item
		return nil
	})
	if err != nil {
		return opname.Item{}, err
	}
	return result, nil
}

func (s Service) GetDetail(ctx context.Context, opnameID string) (opname.Detail, error) {
	header, err := s.Repo.FindOpnameByID(ctx, opnameID)
	if err != nil {
		return opname.Detail{}, apperror.New(http.StatusNotFound, "Data opname tidak ditemukan", err)
	}
	items, err := s.Repo.FindOpnameItems(ctx, opnameID)
	if err != nil {
		return opname.Detail{}, apperror.New(http.StatusInternalServerError, "Gagal mengambil item opname", err)
	}
	return opname.Detail{
		ID:          header.ID,
		Description: header.Description,
		OpnameDate:  header.OpnameDate,
		TotalOpname: header.TotalOpname,
		Items:       items,
	}, nil
}

func (s Service) GetItems(ctx context.Context, opnameID string) ([]opname.Item, error) {
	items, err := s.Repo.FindOpnameItems(ctx, opnameID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Gagal mengambil item opname", err)
	}
	return items, nil
}
