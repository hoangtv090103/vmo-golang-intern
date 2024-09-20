package utils

import (
	"bytes"
	"ecommerce/internal/order/domain"
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePDF(invoiceData domain.InvoiceData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: %d", invoiceData.OrderID))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Order Date: "+invoiceData.OrderDate)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Customer: "+invoiceData.CustomerName)
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(80, 10, "Product Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Unit Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	// Table content
	pdf.SetFont("Arial", "", 12)
	for _, item := range invoiceData.Items {
		pdf.CellFormat(80, 10, item.ProductName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", item.TotalPrice), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetX(-50) // Move the cursor to the right edge minus 50 units
	pdf.CellFormat(40, 10, fmt.Sprintf("Total: %.2f", invoiceData.Total), "1", 0, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
