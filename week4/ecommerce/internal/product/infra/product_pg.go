package infra

import (
	"context"
	"database/sql"
	"ecommerce/internal/product/entity"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductPGRepository struct {
	DB *sql.DB
}

func NewProductPGRepository(db *sql.DB) *ProductPGRepository {
	return &ProductPGRepository{
		DB: db,
	}
}

func (pr *ProductPGRepository) Create(ctx context.Context, product *entity.Product) error {
	if product.Price < 0 {
		return errors.New("invalid price")
	}

	if product.Stock < 0 {
		return errors.New("invalid stock")
	}
	
	query := `INSERT INTO products (name, description, price, stock, image_path) VALUES (?, ?, ?, ?, ?)`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	
	_, err := pr.DB.ExecContext(
		ctx,
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

func (pr *ProductPGRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	query := `SELECT id, COALESCE(name, ''), COALESCE(price, 0.0), COALESCE(stock, 0), COALESCE(description, ''), COALESCE(image_path, '') FROM products`
	rows, err := pr.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var products []*entity.Product

	for rows.Next() {
		product := &entity.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description, &product.ImagePath)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (pr *ProductPGRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {

	product := &entity.Product{}

	query := `SELECT id, COALESCE(name, ''), COALESCE(description, ''), COALESCE(price, 0.0), COALESCE(stock, 0), COALESCE(image_path, '') FROM products WHERE id = $1`

	err := pr.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.ImagePath)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *ProductPGRepository) GetByName(ctx context.Context, name string) ([]*entity.Product, error) {
	var (
		err      error
		products []*entity.Product
		rows     *sql.Rows
	)

	query := `SELECT * FROM products WHERE products.document @@ to_tsquery('?')`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err = pr.DB.QueryContext(
		ctx,
		query,
		name,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product := &entity.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description, &product.ImagePath)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (pr *ProductPGRepository) Update(ctx context.Context, product *entity.Product) error {
	// Initialize the query and arguments
	query := "UPDATE products SET"
	var args []interface{}

	// Check and append non-empty fields
	if product.Name != "" {
		query += " name = ?,"
		args = append(args, product.Name)
	}
	if product.Description != "" {
		query += " description = ?,"
		args = append(args, product.Description)
	}
	if product.Price != 0 {
		query += " price = ?,"
		args = append(args, product.Price)
	}
	if product.Stock != 0 {
		query += " stock = ?,"
		args = append(args, product.Stock)
	}
	if product.ImagePath != "" {
		query += " image_path = ?,"
		args = append(args, product.ImagePath)
	}

	// Remove the trailing comma and add the WHERE clause
	query = strings.TrimSuffix(query, ",")
	query += " WHERE id = ?"
	args = append(args, product.ID)

	// Execute the query
	_, err := pr.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductPGRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := pr.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
