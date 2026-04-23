package member

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/member"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) ExportExcel(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Members.ListMembers(ctx, branchID, member.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export members excel failed", err.Error())
	}
	f := exportshared.NewExcelFile("Members")
	sheet := "Members"
	f.SetCellValue(sheet, "A1", "DATA MEMBERS")
	headers := []string{"ID", "NAME", "PHONE", "ADDRESS", "MEMBER CATEGORY", "POINTS"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	for i, item := range items.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Address)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.MemberCategory)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Points)
	}
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export members excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("members-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID string) ([]byte, string, error) {
	items, err := s.Members.ListMembers(ctx, branchID, member.ListRequest{Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export members pdf failed", err.Error())
	}
	pdf := exportshared.NewPDF("MASTER MEMBERS")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, "MASTER MEMBERS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "NAME", "PHONE", "ADDRESS", "CATEGORY", "POINTS"}
	widths := []float64{35, 50, 35, 80, 45, 32}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 8)
	for _, item := range items.Items {
		values := []string{item.ID, item.Name, item.Phone, item.Address, item.MemberCategory, fmt.Sprintf("%d", item.Points)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export members pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("members-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
