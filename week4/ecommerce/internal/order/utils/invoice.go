package utils

import (
	"bytes"
	"ecommerce/internal/order/entity"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"time"
)

func GenerateInvoicePDF(invoiceData []*entity.InvoiceData) ([]byte, error) {
	invoice := invoiceData[0]
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pageWidth, _ := pdf.GetPageSize()
	pdf.SetX((pageWidth - pdf.GetStringWidth("Invoice")) / 2)
	pdf.CellFormat(pdf.GetStringWidth("Invoice"), 10, "Invoice", "0", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(20)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: %d", invoice.OrderID))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Order Date: "+invoice.OrderDate)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Print Invoice Date: "+time.Now().Format("2006-01-02"))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Customer: "+invoice.CustomerName)
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
	for _, item := range invoice.Items {
		pdf.CellFormat(80, 10, item.ProductName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", item.TotalPrice), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetX(-50) // Move the cursor to the right edge minus 50 units
	pdf.CellFormat(40, 10, fmt.Sprintf("Total: %.2f", invoice.Total), "1", 0, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
