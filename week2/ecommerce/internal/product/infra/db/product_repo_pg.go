package db

import (
	"database/sql"
	"ecommerce/config"
	productDomain "ecommerce/internal/product/domain"
	"strconv"
	"strings"
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
	query := `INSERT INTO products (name, description, price, stock, image_path) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.PG.GetDB().Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImagePath,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepoPG) GetAll() ([]productDomain.Product, error) {
	query := `SELECT id, COALESCE(name, ''), COALESCE(price, 0.0), COALESCE(stock, 0), COALESCE(description, ''), COALESCE(image_path, '') FROM products`
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
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description, &product.ImagePath)

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

	query := `SELECT id, COALESCE(name, ''), COALESCE(description, ''), COALESCE(price, 0.0), COALESCE(stock, 0), COALESCE(image_path, '') FROM products WHERE id = $1`

	err = p.PG.GetDB().QueryRow(
		query,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.ImagePath)

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

	query := `SELECT id, COALESCE(name, ''), COALESCE(price, 0.0), COALESCE(stock, 0), COALESCE(description, ''), COALESCE(image_path, '') FROM products WHERE name ilike '%' || $1 || '%'`

	rows, err = p.PG.GetDB().Query(
		query,
		name,
	)

	if err != nil {
		return []productDomain.Product{}, err
	}

	for rows.Next() {
		var product productDomain.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description, &product.ImagePath)

		if err != nil {
			return []productDomain.Product{}, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepoPG) Update(product productDomain.Product) error {
	// Initialize the query and arguments
	query := "UPDATE products SET"
	args := []interface{}{}
	argCounter := 1

	// Check and append non-empty fields
	if product.Name != "" {
		query += " name = $" + strconv.Itoa(argCounter) + ","
		args = append(args, product.Name)
		argCounter++
	}
	if product.Description != "" {
		query += " description = $" + strconv.Itoa(argCounter) + ","
		args = append(args, product.Description)
		argCounter++
	}
	if product.Price != 0 {
		query += " price = $" + strconv.Itoa(argCounter) + ","
		args = append(args, product.Price)
		argCounter++
	}
	if product.Stock != 0 {
		query += " stock = $" + strconv.Itoa(argCounter) + ","
		args = append(args, product.Stock)
		argCounter++
	}
	if product.ImagePath != "" {
		query += " image_path = $" + strconv.Itoa(argCounter) + ","
		args = append(args, product.ImagePath)
		argCounter++
	}

	// Remove the trailing comma and add the WHERE clause
	query = strings.TrimSuffix(query, ",")
	query += " WHERE id = $" + strconv.Itoa(argCounter)
	args = append(args, product.ID)

	// Execute the query
	_, err := p.PG.GetDB().Exec(query, args...)
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
