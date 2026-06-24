package wallet

import (
	"exampleApp/internal/domain/wallet"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Repository interface {
	GetBalanceByID(uuid string) (float64, error)
	AddOrUpdateWallet(w wallet.WalletRequest) (wallet.Wallet, error)
}
type WalletService struct {
	repo  Repository
	valid *validator.Validate
}

func New(repo Repository) *WalletService {
	return &WalletService{
		repo:  repo,
		valid: validator.New(),
	}
}

func (s *WalletService) GetBalanceByID(uuid string) (float64, error) {
	balance, err := s.repo.GetBalanceByID(uuid)
	if err != nil {
		return 0.0, err
	}
	return balance, nil
}

func (s *WalletService) ChageAmount(req wallet.WalletRequest) (wallet.Wallet, error) {

	if err := s.valid.Struct(req); err != nil {
		return wallet.Wallet{}, fmt.Errorf("Ошибка валидации: %w", err)
	}
	w := wallet.WalletRequest{
		WalletID:      req.WalletID,
		OperationType: req.OperationType,
		Amount:        req.Amount,
	}

	changedWallet, err := s.repo.AddOrUpdateWallet(w)
	if err != nil {
		return wallet.Wallet{}, err
	}
	return changedWallet, nil
}
