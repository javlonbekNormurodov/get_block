package service

import (
	"math/big"

	"eth_tracker/internal/repository"
)

type BalanceService struct {
	repo *repository.BalanceRepository
}

func NewBalanceService(repo *repository.BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) UpdateBalance(address string, value *big.Int) error {
	return s.repo.UpdateBalance(address, value)
}

func (s *BalanceService) GetMostChangedBalance() (string, *big.Int, error) {
	return s.repo.GetMostChangedBalance()
}
