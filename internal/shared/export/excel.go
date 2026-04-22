package export

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

func NewExcelFile(sheet string) *excelize.File {
	f := excelize.NewFile()
	defaultSheet := f.GetSheetName(0)
	if sheet != "" && sheet != defaultSheet {
		f.SetSheetName(defaultSheet, sheet)
	}
	return f
}

func WriteExcel(f *excelize.File) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ExcelColumnName(n int) (string, error) {
	return excelize.ColumnNumberToName(n)
}
