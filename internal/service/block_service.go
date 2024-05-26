package service

import (
	"math/big"

	"eth_tracker/pkg/eth"
)

type BlockService struct {
	ethClient *eth.Client
}

func NewBlockService(ethClient *eth.Client) *BlockService {
	return &BlockService{ethClient: ethClient}
}

func (s *BlockService) GetLatestBlockNumber() (*big.Int, error) {
	return s.ethClient.GetLatestBlockNumber()
}
