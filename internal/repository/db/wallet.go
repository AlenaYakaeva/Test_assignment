package db

import (
	"context"
	"exampleApp/internal/domain/wallet"
	"fmt"
	"time"
)

func (s *Storage) GetBalanceByID(uuid string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var balance float64
	err := s.conn.QueryRow(ctx, "SELECT Amount FROM wallets WHERE WalletID = $1", uuid).Scan(&balance)

	if err != nil {
		return 0.0, fmt.Errorf(ErrNotFoundWallet, err)
	}
	return balance, nil
}

func (s *Storage) AddOrUpdateWallet(w wallet.WalletRequest) (wallet.Wallet, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	koef := 1.0
	if w.OperationType == "WITHDRAW" {
		koef = -1.0
	}

	// Проверяем существование счета по ID.
	balance, err := s.GetBalanceByID(w.WalletID)
	if err != nil {
		// Если счет не существует и операция "списание", то возвращаем ошибку - попытка списания с несуществующего чета
		if koef == -1.0 {
			return wallet.Wallet{}, ErrWithdrawWalletNotFound
		}

		// Если счет не существует и операция "пополнение", то создаем его и записываем сумму.
		var uuid string
		err := s.conn.QueryRow(ctx,
			`INSERT INTO wallets (
		 
		Amount) 
		VALUES ($1) RETURNING WalletID`,
			w.Amount).Scan(&uuid)
		if err != nil {
			return wallet.Wallet{}, fmt.Errorf(ErrNotCreatedWallet, err)
		}
		return wallet.Wallet{
			WalletID: uuid,
			Amount:   w.Amount,
		}, nil
	}

	// Если счет существует и операция "списание", то возвращаем ошибку, если баланс меньше суммы списания
	if koef == -1.0 && balance < w.Amount {
		return wallet.Wallet{}, ErrPaymentFailed
	}

	// Если счет существуе и операция "списание", то уменьшаем баланс на сумму списания
	// Если счет существует и операция "пополнение", то увеличиваем баланс на сумму пополнения
	w.Amount = koef * w.Amount

	_, err = s.conn.Exec(ctx, `UPDATE wallets 
	SET Amount=Amount+$1 
	WHERE WalletID=$2`, w.Amount, w.WalletID)
	if err != nil {
		return wallet.Wallet{}, fmt.Errorf(ErrNotUpdatedWallet, err)
	}
	return wallet.Wallet{
		WalletID: w.WalletID,
		Amount:   balance + w.Amount,
	}, nil
}
