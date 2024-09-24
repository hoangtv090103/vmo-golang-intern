package reports

import (
	"ecommerce/config"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
)

type ReportRepoPG struct {
	PG *config.PG
}

func NewReportRepoPG(pg *config.PG) *ReportRepoPG {
	return &ReportRepoPG{
		PG: pg,
	}
}

func (r *ReportRepoPG) GenerateReport(startDate, endDate string) (string, error) {
	var e ExcelReport

	// Create necessary directories
	outputDir := "internal/reports"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Generate Excel report
	f := excelize.NewFile()

	// Create a new sheet.
	sheet := "Sheet1"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return "", err
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Merge cells from A1 to D1
	if err := f.MergeCell("Sheet1", "A1", "D1"); err != nil {
		return "", err
	}

	query := `SELECT u.name,
				   count(o.id)        as total_orders,
				   avg(o.total_price) as avg_order_value,
				   sum(o.total_price) as total_spent
			   FROM users u
			       LEFT JOIN orders o ON u.id = o.user_id
			   WHERE o.created_at BETWEEN $1 AND $2
			   GROUP BY u.name
			   ORDER BY total_spent DESC NULLS LAST;`

	// Get Report Date
	rows, err := r.PG.GetDB().Query(query, startDate, endDate)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Set value of a cell.
	f.SetCellValue(sheet, "A1", "Period Report")
	f.SetCellValue(sheet, "A2", "From")
	f.SetCellValue(sheet, "B2", startDate)
	f.SetCellValue(sheet, "A3", "As at")
	f.SetCellValue(sheet, "B3", endDate)

	// Header
	f.SetCellValue(sheet, "A5", "User")
	f.SetCellValue(sheet, "B5", "Number of Order")
	f.SetCellValue(sheet, "C5", "Average Order Value")
	f.SetCellValue(sheet, "D5", "Total Spent")

	// Body
	line := 6
	for rows.Next() {
		err := rows.Scan(&e.Name, &e.TotalOrders, &e.AverageOrderValue, &e.TotalSpent)
		if err != nil {
			return "", err
		}

		name := ""
		if e.Name.Valid {
			name = e.Name.String
		}

		totalOrders := int64(0)
		if e.TotalOrders.Valid {
			totalOrders = e.TotalOrders.Int64
		}

		avgOrderValue := float64(0)
		if e.AverageOrderValue.Valid {
			avgOrderValue = e.AverageOrderValue.Float64
		}

		totalSpent := float64(0)
		if e.TotalSpent.Valid {
			totalSpent = e.TotalSpent.Float64
		}

		f.SetCellValue(sheet, "A"+strconv.Itoa(line), name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(line), totalOrders)
		f.SetCellValue(sheet, "C"+strconv.Itoa(line), avgOrderValue)
		f.SetCellValue(sheet, "D"+strconv.Itoa(line), totalSpent)
		line++
	}

	// Save xlsx file by the given path.
	if err := f.SaveAs("internal/reports/report.xlsx"); err != nil {
		return "", err
	}

	log.Println("Excel report generated successfully")
	// Return the path to the generated report
	return "internal/reports/report.xlsx", nil
}
