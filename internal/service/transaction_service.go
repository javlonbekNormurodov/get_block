package service

import (
	"math/big"

	"eth_tracker/internal/repository"
	"eth_tracker/pkg/eth"
)

type TransactionService struct {
	repo           *repository.TransactionRepository
	ethClient      *eth.Client
	balanceService *BalanceService
}

func NewTransactionService(repo *repository.TransactionRepository, ethClient *eth.Client, balanceService *BalanceService) *TransactionService {
	return &TransactionService{repo: repo, ethClient: ethClient, balanceService: balanceService}
}

func (s *TransactionService) ProcessBlock(blockNumber int64) error {
	block, err := s.ethClient.GetBlockByNumber(blockNumber)
	if err != nil {
		return err
	}

	for _, txHash := range block.Txns {
		tx, err := s.ethClient.GetTransactionByHash(txHash)
		if err != nil {
			return err
		}

		value := new(big.Int)
		value.SetString(tx.Value[2:], 16)

		s.balanceService.UpdateBalance(tx.From, value.Neg(value))
		s.balanceService.UpdateBalance(tx.To, value)

		err = s.repo.SaveTransaction(*tx, blockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}
