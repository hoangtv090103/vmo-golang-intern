package infra

import (
	"context"
	"database/sql"
	"ecommerce/internal/order/entity"
	"errors"
)

type OrderPGRepository struct {
	DB *sql.DB
}

func NewOrderPGRepository(db *sql.DB) *OrderPGRepository {
	return &OrderPGRepository{
		DB: db,
	}
}

func (r *OrderPGRepository) LockProductForUpdate(ctx context.Context, tx *sql.Tx, id int) error {
	query := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`

	var stock int
	var price float64
	err := tx.QueryRowContext(ctx, query, id).Scan(&stock, &price)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPGRepository) BuyProduct(ctx context.Context, tx *sql.Tx, productID, buyQty int) error {
	var stock int
	err := tx.QueryRowContext(ctx, `SELECT stock FROM products WHERE id = $1`, productID).Scan(&stock)
	if err != nil {
		return err
	}

	if stock < buyQty {
		return errors.New("insufficient stock")
	}

	_, err = tx.Exec(`UPDATE products SET stock = stock - $1 WHERE id = $2`, buyQty, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderPGRepository) UpdateUserBalance(ctx context.Context, tx *sql.Tx, userID int, totalPrice float64) error {
	var balance float64
	err := tx.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&balance)
	if err != nil {
		return err
	}

	balance -= totalPrice
	if balance < 0 {
		return errors.New("insufficient balance")
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = $1 WHERE id = $2`, balance, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPGRepository) CreateOrderLine(ctx context.Context, tx *sql.Tx, orderID int, line entity.OrderLine) error {
	query := `INSERT INTO order_lines (order_id, product_id, qty, total) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := tx.ExecContext(ctx, query, orderID, line.ProductID, line.Qty, line.Total)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPGRepository) Create(ctx context.Context, order *entity.Order) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRowContext(ctx, `INSERT INTO orders (user_id) VALUES ($1) RETURNING id, created_at`, order.UserID).Scan(&order.ID, &order.OrderDate)
	if err != nil {
		return err
	}

	totalPrice := 0.0
	for _, line := range order.Lines {
		err = tx.QueryRowContext(ctx, `SELECT price FROM products WHERE id = $1`, line.ProductID).Scan(&line.Product.Price)
		if err != nil {
			return err
		}

		line.Total = line.Product.Price * float64(line.Qty)
		totalPrice += line.Total

		err = r.BuyProduct(ctx, tx, line.ProductID, line.Qty)
		if err != nil {
			return err
		}

		err = r.CreateOrderLine(ctx, tx, order.ID, line)
		if err != nil {
			return err
		}
	}

	err = r.UpdateUserBalance(ctx, tx, order.UserID, totalPrice)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE orders SET total_price = $1 WHERE id = $2`, totalPrice, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPGRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	query := `SELECT id, user_id, created_at, total_price FROM orders`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalPrice)
		if err != nil {
			return nil, err
		}

		order.Lines, err = r.getOrderLines(ctx, order.ID)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderPGRepository) getOrderLines(ctx context.Context, orderID int) ([]entity.OrderLine, error) {
	query := `SELECT id, order_id, product_id, qty, total FROM order_lines WHERE order_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []entity.OrderLine
	for rows.Next() {
		line := entity.OrderLine{}
		err := rows.Scan(&line.ID, &line.OrderID, &line.ProductID, &line.Qty, &line.Total)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (r *OrderPGRepository) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	query := `SELECT id, user_id, created_at, total_price FROM orders WHERE id = $1`
	order := &entity.Order{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	order.Lines, err = r.getOrderLines(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderPGRepository) GetUserOrders(ctx context.Context, username string) ([]*entity.Order, error) {
	query := `SELECT o.id, o.user_id, o.created_at, o.total_price FROM orders o JOIN users u ON o.user_id = u.id WHERE u.username = $1`
	rows, err := r.DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalPrice)
		if err != nil {
			return nil, err
		}

		order.Lines, err = r.getOrderLines(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderPGRepository) Update(ctx context.Context, order *entity.Order) error {
	query := `UPDATE orders SET user_id = $1, total_price = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, order.UserID, order.TotalPrice, order.ID)
	if err != nil {
		return err
	}

	for _, line := range order.Lines {
		err = r.updateOrderLine(ctx, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderPGRepository) updateOrderLine(ctx context.Context, line entity.OrderLine) error {
	query := `UPDATE order_lines SET product_id = $1, qty = $2, total = $3 WHERE id = $4`
	_, err := r.DB.ExecContext(ctx, query, line.ProductID, line.Qty, line.Total, line.ID)
	return err
}

func (r *OrderPGRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *OrderPGRepository) GetInvoice(ctx context.Context, orderID int) ([]*entity.InvoiceData, error) {
	query := `SELECT o.id, o.created_at, u.username, ol.product_id, ol.qty, ol.total, p.name, p.price
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN order_lines ol ON o.id = ol.order_id
		JOIN products p ON ol.product_id = p.id
		WHERE o.id = $1`

	rows, err := r.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make(map[int]*entity.InvoiceData)
	for rows.Next() {
		var (
			orderID      int
			orderDate    string
			customerName string
			productID    int
			qty          int
			total        float64
			productName  string
			unitPrice    float64
		)
		err := rows.Scan(&orderID, &orderDate, &customerName, &productID, &qty, &total, &productName, &unitPrice)
		if err != nil {
			return nil, err
		}

		if _, exists := invoices[orderID]; !exists {
			invoices[orderID] = &entity.InvoiceData{
				OrderID:      orderID,
				OrderDate:    orderDate,
				CustomerName: customerName,
				Items:        []entity.InvoiceItem{},
				Total:        0,
			}
		}

		invoice := invoices[orderID]
		invoice.Items = append(invoice.Items, entity.InvoiceItem{
			ProductName: productName,
			Quantity:    qty,
			UnitPrice:   unitPrice,
			TotalPrice:  total,
		})
		invoice.Total += total
	}

	var result []*entity.InvoiceData
	for _, invoice := range invoices {
		result = append(result, invoice)
	}

	return result, nil
}