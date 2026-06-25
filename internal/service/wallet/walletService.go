package wallet

import (
	"context"
	"exampleApp/internal/domain/wallet"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Repository interface {
	GetBalanceByID(uuid string) (float64, error)
	AddOrUpdateWallet(ctx context.Context, w wallet.WalletRequest) (wallet.Wallet, error)
}
type depositRequest struct {
	walletReq wallet.WalletRequest
	resWallet chan wallet.Wallet
	response  chan error // Канал обратной связи. Через него воркер передаст результат (успех/ошибка) обратно в HTTP-поток.
}
type WalletService struct {
	repo    Repository
	mu      sync.RWMutex
	wallets map[string]chan depositRequest
	valid   *validator.Validate
}

func New(repo Repository) *WalletService {
	return &WalletService{
		repo:    repo,
		wallets: make(map[string]chan depositRequest),
		valid:   validator.New(),
	}
}

func (s *WalletService) GetBalanceByID(uuid string) (float64, error) {
	balance, err := s.repo.GetBalanceByID(uuid)
	if err != nil {
		return 0.0, err
	}
	return balance, nil
}

func (s *WalletService) ChageAmount(ctx context.Context, walReq wallet.WalletRequest) (wallet.Wallet, error) {
	if err := s.valid.Struct(walReq); err != nil {
		return wallet.Wallet{}, fmt.Errorf("Ошибка валидации: %w", err)
	}

	s.mu.Lock()
	ch, exists := s.wallets[walReq.WalletID]
	if !exists {
		ch = make(chan depositRequest, 5000)
		s.wallets[walReq.WalletID] = ch
		go s.startWalletWorker(ch)
	}
	s.mu.Unlock()

	resCh := make(chan error, 1)
	resWalCh := make(chan wallet.Wallet, 1)
	req := depositRequest{
		walletReq: walReq,
		resWallet: resWalCh,
		response:  resCh,
	}
	select {
	case ch <- req:
	case <-ctx.Done():
		return wallet.Wallet{}, ctx.Err()
	}

	select {
	case wal := <-resWalCh:
		err := <-resCh
		return wal, err
	case <-ctx.Done():
		return wallet.Wallet{}, ctx.Err()
	}
}

func (s *WalletService) startWalletWorker(ch chan depositRequest) {

	for req := range ch {

		w, err := s.repo.AddOrUpdateWallet(context.Background(), req.walletReq)

		req.resWallet <- w
		req.response <- err
	}
}
