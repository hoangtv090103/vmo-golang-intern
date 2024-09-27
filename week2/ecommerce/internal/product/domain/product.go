package domain

// Product struct represents the product entity
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImagePath   string  `json:"image_path"`
}

func (p *Product) GetID() int {
	return p.ID
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) GetDescription() string {
	return p.Description
}

func (p *Product) GetPrice() float64 {
	return p.Price
}

func (p *Product) GetStock() int {
	return p.Stock
}

func (p *Product) GetImagePath() string { return p.ImagePath }

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) SetDescription(description string) {
	p.Description = description
}

func (p *Product) SetPrice(price float64) {
	p.Price = price
}

func (p *Product) SetStock(stock int) {
	p.Stock = stock
}

func (p *Product) SetImagePath(path string) {
	p.ImagePath = path
}
