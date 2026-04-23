package member

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/member"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Members ports.MemberRepository
	IDs     ports.IDGenerator
}

func (s Service) List(ctx context.Context, branchID string, req member.ListRequest) (member.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Members.ListMembers(ctx, branchID, req)
	if err != nil {
		return member.ListResult{}, apperror.New(http.StatusInternalServerError, "Get members failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (member.Member, error) {
	item, err := s.Members.FindMemberDetailByID(ctx, id, branchID)
	if err != nil {
		return member.Member{}, apperror.New(http.StatusNotFound, "Get member failed", "member not found")
	}
	return item, nil
}

func (s Service) Create(ctx context.Context, branchID string, req member.CreateRequest) (member.Member, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return member.Member{}, apperror.New(http.StatusBadRequest, "Create member failed", "name is required")
	}
	if req.MemberCategoryID == 0 {
		return member.Member{}, apperror.New(http.StatusBadRequest, "Create member failed", "member_category_id is required")
	}
	item := member.Member{
		ID:               s.IDs.New("MBR"),
		Name:             name,
		Phone:            strings.TrimSpace(req.Phone),
		Address:          strings.TrimSpace(req.Address),
		MemberCategoryID: strconv.FormatUint(uint64(req.MemberCategoryID), 10),
		Points:           req.Points,
		BranchID:         branchID,
	}
	if err := s.Members.CreateMember(ctx, item); err != nil {
		return member.Member{}, apperror.New(http.StatusInternalServerError, "Create member failed", err.Error())
	}
	return s.Members.FindMemberDetailByID(ctx, item.ID, branchID)
}

func (s Service) Update(ctx context.Context, branchID, id string, req member.CreateRequest) (member.Member, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return member.Member{}, apperror.New(http.StatusBadRequest, "Update member failed", "name is required")
	}
	if req.MemberCategoryID == 0 {
		return member.Member{}, apperror.New(http.StatusBadRequest, "Update member failed", "member_category_id is required")
	}
	if _, err := s.Members.FindMemberDetailByID(ctx, id, branchID); err != nil {
		return member.Member{}, apperror.New(http.StatusNotFound, "Update member failed", "member not found")
	}
	item := member.Member{ID: id, Name: name, Phone: strings.TrimSpace(req.Phone), Address: strings.TrimSpace(req.Address), MemberCategoryID: strconv.FormatUint(uint64(req.MemberCategoryID), 10), Points: req.Points, BranchID: branchID}
	if err := s.Members.UpdateMember(ctx, item); err != nil {
		return member.Member{}, apperror.New(http.StatusInternalServerError, "Update member failed", err.Error())
	}
	return s.Members.FindMemberDetailByID(ctx, id, branchID)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	if _, err := s.Members.FindMemberDetailByID(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusNotFound, "Delete member failed", "member not found")
	}
	if err := s.Members.DeleteMember(ctx, id, branchID); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete member failed", err.Error())
	}
	return nil
}

func (s Service) Combo(ctx context.Context, branchID, search string) ([]member.ComboItem, error) {
	items, err := s.Members.GetMemberCombo(ctx, branchID, search)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get member combo failed", err.Error())
	}
	return items, nil
}
