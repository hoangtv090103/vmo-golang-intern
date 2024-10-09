package entity

type InvoiceData struct {
	OrderID      int
	OrderDate    string
	CustomerName string
	Items        []InvoiceItem
	Total        float64
}
