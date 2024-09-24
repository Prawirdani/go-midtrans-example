package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/prawirdani/go-midtrans-example/internal/entity"
)

type TransactionRepository interface {
	Insert(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error)
	SelectByID(ctx context.Context, id string) (*entity.Transaction, error)
	Select(ctx context.Context) ([]entity.Transaction, error)
	SaveChanges(ctx context.Context, transaction entity.Transaction) error
}

type transactionRepository struct {
	conn        *sql.DB
	productRepo ProductRepository
}

func NewTransactionRepository(conn *sql.DB, productRepo ProductRepository) TransactionRepository {
	return &transactionRepository{
		conn:        conn,
		productRepo: productRepo,
	}
}

// TODO: Should optimize this mess, probably by running some db call into single call
func (r *transactionRepository) Insert(
	ctx context.Context,
	transaction entity.Transaction,
) (*entity.Transaction, error) {
	err := useTX(r.conn, func(tx *sql.Tx) error {
		// Insert transaction
		tQuery := "INSERT INTO transactions (id, user_id, total) VALUES (?, ?, ?)"
		if _, err := tx.ExecContext(ctx, tQuery, transaction.ID, transaction.User.ID, transaction.Total); err != nil {
			log.Println("Insert.Transaction", err)
			return err
		}

		// Bulk insert transaction details
		tdQuery := "INSERT INTO transaction_details (transaction_id, product_id, product_price, quantity, subtotal) VALUES "
		tdArgs := []interface{}{}
		for i, detail := range transaction.Details {
			if i > 0 {
				tdQuery += ", "
			}
			tdQuery += "(?, ?, ?, ?, ?)"
			tdArgs = append(
				tdArgs,
				transaction.ID,
				detail.Product.ID,
				detail.Product.Price,
				detail.Quantity,
				detail.Subtotal,
			)
		}

		_, err := tx.ExecContext(ctx, tdQuery, tdArgs...)
		if err != nil {
			log.Println("Insert.TransactionDetails", err)
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.SelectByID(ctx, transaction.ID.String())
}

func (r *transactionRepository) SelectByID(
	ctx context.Context,
	id string,
) (*entity.Transaction, error) {
	query := querySelectTransaction + `
	WHERE t.id = ?
	GROUP BY t.id, u.id;
	`

	var transaction entity.Transaction
	var detailJSON string
	err := r.conn.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.Total,
		&transaction.Status,
		&transaction.User.ID,
		&transaction.User.FirstName,
		&transaction.User.LastName,
		&transaction.User.Email,
		&transaction.User.Phone,
		&detailJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrTransactionNotFound
		}
		return nil, err
	}

	details := make([]entity.TransactionDetails, 0)
	if err := json.Unmarshal([]byte(detailJSON), &details); err != nil {
		return nil, err
	}
	transaction.Details = details

	return &transaction, nil
}

func (r *transactionRepository) Select(ctx context.Context) ([]entity.Transaction, error) {
	transactions := make([]entity.Transaction, 0)

	query := querySelectTransaction + `GROUP BY t.id, u.id;`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		log.Println("List.Transaction", err)
	}

	for rows.Next() {
		var transaction entity.Transaction
		var detailJSON string

		err := rows.Scan(
			&transaction.ID,
			&transaction.Total,
			&transaction.Status,
			&transaction.User.ID,
			&transaction.User.FirstName,
			&transaction.User.LastName,
			&transaction.User.Email,
			&transaction.User.Phone,
			&detailJSON,
		)
		if err != nil {
			log.Println("Scan.Transaction", err)
			return nil, err
		}

		details := make([]entity.TransactionDetails, 0)
		if err := json.Unmarshal([]byte(detailJSON), &details); err != nil {
			log.Println("Unmarshal.Transaction", err)
			return nil, err
		}

		transaction.Details = details
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// SaveChanges() is used to update transaction status
func (r *transactionRepository) SaveChanges(
	ctx context.Context,
	transaction entity.Transaction,
) error {
	query := "UPDATE transactions SET status_id=(SELECT id from transaction_status WHERE status=?) WHERE id = ?"
	_, err := r.conn.ExecContext(ctx, query, transaction.Status, transaction.ID)
	if err != nil {
		log.Println("Update.TransactionStatus", err)
		return err
	}

	return nil
}

const querySelectTransaction = `
SELECT 
	t.id, 
	t.total,
	ts.status,
	u.id AS user_id, 
	u.first_name, 
	u.last_name, 
	u.email, 
	u.phone,
	json_group_array(
		json_object(
			'id', td.id,
			'quantity', td.quantity,
			'subtotal', td.subtotal,
			'product', json_object(
				'id', td.product_id,
				'price', td.product_price,
				'name', p.name
			)
		)
	) as details
FROM transactions t
JOIN transaction_status ts ON t.status_id = ts.id
JOIN users u ON t.user_id = u.id
JOIN transaction_details td ON t.id = td.transaction_id
JOIN products p ON td.product_id = p.id
`
