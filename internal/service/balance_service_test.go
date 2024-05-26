package service

import (
	"math/big"
	"testing"

	"eth_tracker/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBalanceService_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	balanceRepo := repository.NewBalanceRepository(db)
	balanceService := NewBalanceService(balanceRepo)

	address := "0x123"
	value := big.NewInt(100)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance FROM balances WHERE address = ?").
		WithArgs(address).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("0"))
	mock.ExpectExec("UPDATE balances SET balance = ? WHERE address = ?").
		WithArgs(value.String(), address).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = balanceService.UpdateBalance(address, value)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBalanceService_GetMostChangedBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	balanceRepo := repository.NewBalanceRepository(db)
	balanceService := NewBalanceService(balanceRepo)

	address := "0x123"
	totalChange := big.NewInt(200)

	mock.ExpectQuery("SELECT address, SUM(ABS(value)) as total_change FROM transactions GROUP BY address ORDER BY total_change DESC LIMIT 1").
		WillReturnRows(sqlmock.NewRows([]string{"address", "total_change"}).AddRow(address, totalChange.String()))

	addr, change, err := balanceService.GetMostChangedBalance()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if addr != address {
		t.Errorf("expected address %s, got %s", address, addr)
	}

	if change.Cmp(totalChange) != 0 {
		t.Errorf("expected change %s, got %s", totalChange.String(), change.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
