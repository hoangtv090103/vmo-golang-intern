package entity

// Product struct represents the product entity
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImagePath   string  `json:"image_path"`
}

// NewProduct creates a new product entity
func NewProduct() *Product {
	return &Product{
		ID:          0,
		Name:        "",
		Description: "",
		Price:       0,
		Stock:       0,
		ImagePath:   "",
	}
}

func (p *Product) SetName(name string) *Product {
	p.Name = name
	return p
}

func (p *Product) SetDescription(description string) *Product {
	p.Description = description
	return p
}

func (p *Product) SetPrice(price float64) *Product {
	p.Price = price
	return p
}

func (p *Product) SetStock(stock int) *Product {
	p.Stock = stock
	return p
}

func (p *Product) SetImagePath(path string) *Product {
	p.ImagePath = path
	return p
}
