package domain

type InvoiceData struct {
	OrderID      int
	OrderDate    string
	CustomerName string
	Items        []InvoiceItem
	Total        float64
}

type InvoiceItem struct {
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}
