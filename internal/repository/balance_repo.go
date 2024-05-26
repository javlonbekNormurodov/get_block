package repository

import (
	"database/sql"
	"math/big"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) UpdateBalance(address string, value *big.Int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var balance big.Int
	err = tx.QueryRow("SELECT balance FROM balances WHERE address = $1", address).Scan(&balance)
	if err == sql.ErrNoRows {
		_, err = tx.Exec("INSERT INTO balances (address, balance) VALUES ($1, $2)", address, value)
	} else if err == nil {
		newBalance := new(big.Int).Add(&balance, value)
		_, err = tx.Exec("UPDATE balances SET balance = $1 WHERE address = $2", newBalance, address)
	}

	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}

	return nil
}

func (r *BalanceRepository) GetMostChangedBalance() (string, *big.Int, error) {
	rows, err := r.db.Query(`
        SELECT address, SUM(ABS(value)) as total_change
        FROM transactions
        GROUP BY address
        ORDER BY total_change DESC
        LIMIT 1
    `)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var address string
	var totalChange big.Int
	if rows.Next() {
		err := rows.Scan(&address, &totalChange)
		if err != nil {
			return "", nil, err
		}
	}

	return address, &totalChange, nil
}
