package db

import (
	"database/sql"
	"ecommerce/config"
	productDomain "ecommerce/internal/product/domain"
	// userDomain "ecommerce/internal/user/domain"
	// "errors"
)

type ProductRepoPG struct {
	PG *config.PG
}

func NewProductRepoPG(pg *config.PG) *ProductRepoPG {
	return &ProductRepoPG{
		PG: pg,
	}
}

func (p *ProductRepoPG) Create(product productDomain.Product) error {
	query := `INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)`
	_, err := p.PG.GetDB().Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepoPG) GetAll() ([]productDomain.Product, error) {
	query := `SELECT id, name, price, stock, description FROM products`
	rows, err := p.PG.GetDB().Query(query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var products []productDomain.Product

	for rows.Next() {
		var product productDomain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepoPG) GetByID(id int) (productDomain.Product, error) {
	var (
		err     error
		product productDomain.Product
	)

	query := `SELECT id, name, description, price, stock FROM products WHERE id = $1`

	err = p.PG.GetDB().QueryRow(
		query,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)

	if err != nil {
		return productDomain.Product{}, err
	}

	return product, nil
}

func (p *ProductRepoPG) GetByName(name string) ([]productDomain.Product, error) {
	var (
		err      error
		products []productDomain.Product
		rows     *sql.Rows
	)

	query := `SELECT id, name, description, price, stock from products WHERE name ilike '%' || $1 || '%'`

	rows, err = p.PG.GetDB().Query(
		query,
		name,
	)

	if err != nil {
		return []productDomain.Product{}, err
	}

	for rows.Next() {
		var product productDomain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)

		if err != nil {
			return []productDomain.Product{}, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepoPG) Update(product productDomain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, stock = $4`

	_, err := p.PG.GetDB().Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepoPG) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := p.PG.GetDB().Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
