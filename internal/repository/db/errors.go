package db

import "errors"

var (
	ErrWithdrawWalletNotFound = errors.New("Невозможно списать деньги с несуществующего счета")
	ErrPaymentFailed          = errors.New("Невозможно провести списание")
)

const (
	ErrNotCreatedWallet = "Невозможно создать счет %w"
	ErrNotUpdatedWallet = "Невозможно обновить счет %w"
	ErrNotFoundWallet   = "Счет не найден: %w"
)
