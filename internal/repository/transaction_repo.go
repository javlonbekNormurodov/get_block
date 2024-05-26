package repository

import (
	"database/sql"
	"eth_tracker/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) SaveTransaction(tx domain.Transaction, blockNumber int64) error {
	_, err := r.db.Exec(`
        INSERT INTO transactions (block_number, address, value)
        VALUES ($1, $2, $3)
        ON CONFLICT (block_number, address) DO UPDATE
        SET value = transactions.value + EXCLUDED.value
    `, blockNumber, tx.From, tx.Value)
	return err
}
