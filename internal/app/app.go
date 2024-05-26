package app

import (
	"database/sql"
	"eth_tracker/internal/repository"
	"eth_tracker/internal/service"
	"eth_tracker/pkg/eth"
	"log"
	"math/big"
)

type Application struct {
	BalanceService     *service.BalanceService
	BlockService       *service.BlockService
	TransactionService *service.TransactionService
}

func NewApplication(db *sql.DB, ethClient *eth.Client) *Application {
	balanceRepo := repository.NewBalanceRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	balanceService := service.NewBalanceService(balanceRepo)
	blockService := service.NewBlockService(ethClient)
	transactionService := service.NewTransactionService(transactionRepo, ethClient, balanceService)

	return &Application{
		BalanceService:     balanceService,
		BlockService:       blockService,
		TransactionService: transactionService,
	}
}

func (app *Application) Run() error {
	latestBlockNumber, err := app.BlockService.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	startBlock := new(big.Int).Sub(latestBlockNumber, big.NewInt(100))

	for i := startBlock.Int64(); i <= latestBlockNumber.Int64(); i++ {
		err := app.TransactionService.ProcessBlock(i)
		if err != nil {
			log.Printf("Error processing block %d: %v", i, err)
		}
	}

	address, change, err := app.BalanceService.GetMostChangedBalance()
	if err != nil {
		return err
	}

	log.Printf("Address with the most balance change: %s, Change: %s", address, change.String())
	return nil
}
