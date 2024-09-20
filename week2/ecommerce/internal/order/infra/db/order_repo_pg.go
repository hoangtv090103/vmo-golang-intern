package db

import (
	"database/sql"
	"ecommerce/config"
	"ecommerce/internal/order/domain"
	"errors"
)

type OrderRepoPG struct {
	PG *config.PG
}

func NewOrderRepoPG(pg *config.PG) *OrderRepoPG {
	return &OrderRepoPG{
		PG: pg,
	}
}

func (o *OrderRepoPG) LockProductForUpdate(tx *sql.Tx, id int) error {
	query := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`

	var stock int
	var price float64
	err := tx.QueryRow(query, id).Scan(&stock, &price)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepoPG) BuyProduct(tx *sql.Tx, productID, buyQty int) error {
	// Get product stock
	var stock int
	err := tx.QueryRow(`SELECT stock FROM products WHERE id = $1`, productID).Scan(&stock)
	if err != nil {
		return err
	}

	// Check if stock is sufficient
	if stock < buyQty {
		return errors.New("insufficient stock")
	}

	_, err = tx.Exec(`UPDATE products SET stock = stock - $1 WHERE id = $2`, buyQty, productID)

	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepoPG) UpdateUserBalance(tx *sql.Tx, userID int, totalPrice float64) error {
	// Lock user and update balance
	var balance float64
	var err error

	err = tx.QueryRow(`SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&balance)

	// Update user balance after order placement
	balance -= totalPrice
	if balance < 0 {
		return errors.New("insufficient balance")
	}

	_, err = tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2`, balance, userID)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepoPG) CreateOrderLine(tx *sql.Tx, orderID int, line domain.OrderLine) error {
	query := `INSERT INTO order_lines (order_id, product_id, qty, total) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := tx.Exec(query, orderID, line.ProductID, line.Qty, line.Total)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepoPG) Create(order domain.Order) error {
	// Begin a transaction
	tx, err := o.PG.GetDB().Begin()
	if err != nil {
		return err
	}
	// Ensure we commit or rollback the transaction properly
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err := tx.Commit()
		if err != nil {
			return
		}
	}()

	// Insert the order
	err = tx.QueryRow(`INSERT INTO orders (user_id) VALUES ($1) RETURNING id, created_at`,
		order.UserID).Scan(&order.ID, &order.OrderDate)
	if err != nil {
		return err
	}

	// Calculate total price and update product stock
	totalPrice := 0.0

	for _, line := range order.Lines {
		// Lock product and retrieve stock/price
		//err = o.LockProductForUpdate(tx, line.ProductID)
		//if err != nil {
		//	return err
		//}

		// Retrieve product
		err = tx.QueryRow(`SELECT price FROM products WHERE id = $1`, line.ProductID).Scan(&line.Product.Price)
		if err != nil {
			return err
		}

		// Update stock and calculate total line price
		line.Total = line.Product.Price * float64(line.Qty)
		totalPrice += line.Total

		// Update product stock
		err = o.BuyProduct(tx, line.ProductID, line.Qty)
		if err != nil {
			return err
		}

		// Insert order line
		err = o.CreateOrderLine(tx, order.ID, line)
		if err != nil {
			return err
		}
	}

	// Update user balance
	err = o.UpdateUserBalance(tx, order.UserID, totalPrice)
	if err != nil {
		return err
	}

	// Update order total price
	_, err = tx.Exec(`UPDATE orders SET total_price = $1 WHERE id = $2`, totalPrice, order.ID)
	if err != nil {
		return err
	}

	return nil
}
func (o *OrderRepoPG) GetAll() ([]domain.Order, error) {
	query := `SELECT id, user_id, created_at, total_price FROM orders`
	queryLines := `SELECT id, product_id, qty, total FROM order_lines WHERE order_id = $1`
	rows, err := o.PG.GetDB().Query(query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var orders []domain.Order

	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalPrice)
		if err != nil {
			return nil, err
		}

		// Retrieve user details
		err = o.PG.GetDB().QueryRow(`SELECT id, name, username, email, balance FROM users WHERE id = $1`, order.UserID).Scan(
			&order.User.ID, &order.User.Name, &order.User.Username, &order.User.Email, &order.User.Balance)
		if err != nil {
			return nil, err
		}

		// Retrieve order lines
		rowsLines, err := o.PG.GetDB().Query(queryLines, order.ID)
		if err != nil {
			return nil, err
		}
		defer func(rowsLines *sql.Rows) {
			err := rowsLines.Close()
			if err != nil {
				return
			}
		}(rowsLines)

		for rowsLines.Next() {
			var line domain.OrderLine
			err := rowsLines.Scan(&line.ID, &line.ProductID, &line.Qty, &line.Total)
			if err != nil {
				return nil, err
			}

			// Retrieve product details
			err = o.PG.GetDB().QueryRow(`SELECT id, name, price, stock FROM products WHERE id = $1`, line.ProductID).Scan(
				&line.Product.ID, &line.Product.Name, &line.Product.Price, &line.Product.Stock)
			if err != nil {
				return nil, err
			}

			order.Lines = append(order.Lines, line)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderRepoPG) GetByID(id int) (domain.Order, error) {
	var order domain.Order
	var orderLines []domain.OrderLine
	// Retrieve order details
	err := o.PG.GetDB().QueryRow(
		`SELECT id, user_id, created_at, total_price FROM orders WHERE id = $1`,
		id,
	).Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalPrice)

	if err != nil {
		return domain.Order{}, err
	}

	// Retrieve lines for the order
	rows, err := o.PG.GetDB().Query(
		`SELECT id, product_id, qty, total FROM order_lines WHERE order_id = $1`,
		id,
	)

	if err != nil {
		return domain.Order{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var line domain.OrderLine
		err := rows.Scan(&line.ID, &line.ProductID, &line.Qty, &line.Total)
		if err != nil {
			return domain.Order{}, err
		}

		orderLines = append(orderLines, line)
	}

	order.Lines = orderLines

	return order, nil
}

func (o *OrderRepoPG) Update(order domain.Order) error {
	query := `UPDATE orders SET user_id = $1, total = $2 WHERE id = $3`

	_, err := o.PG.GetDB().Exec(
		query,
		order.UserID,
		order.TotalPrice,
		order.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepoPG) Delete(id int) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := o.PG.GetDB().Exec(
		query,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
