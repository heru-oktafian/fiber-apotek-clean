package export

import (
	"bytes"

	"github.com/jung-kurt/gofpdf"
)

func NewPDF(title string) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle(title, false)
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()
	return pdf
}

func WritePDF(pdf *gofpdf.Fpdf) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
