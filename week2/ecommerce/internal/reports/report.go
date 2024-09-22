package reports

import (
	"database/sql"
)

type ExcelReport struct {
	Name              sql.NullString
	TotalOrders       sql.NullInt64
	AverageOrderValue sql.NullFloat64
	TotalSpent        sql.NullFloat64
}

type ReportRepository interface {
	GenerateReport(startDate, endDate string) (string, error)
}
