package usecase

import (
	"context"
	"fmt"

	"wallet-service/internal/model"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

type WalletUseCase struct {
	repo repository.WalletRepository
}

func NewWalletUseCase(repo repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

func (u *WalletUseCase) ProcessOperation(
	ctx context.Context,
	walletID uuid.UUID,
	opType model.OperationType,
	amountCents int64,
) (*model.OperationResult, error) {
	if amountCents <= 0 {
		return nil, model.ErrInvalidAmount
	}
	if opType != model.OperationDeposit && opType != model.OperationWithdraw {
		return nil, model.ErrInvalidOperationType
	}

	txRepo, err := u.repo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		txRepo.Rollback(ctx)
	}()

	wallet, err := txRepo.GetWalletForUpdate(ctx, walletID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			wallet = &model.Wallet{
				WalletID: walletID,
				Balance:  0,
			}
			if err := txRepo.CreateWallet(ctx, wallet); err != nil {
				return nil, fmt.Errorf("create wallet: %w", err)
			}
		} else {
			return nil, fmt.Errorf("get wallet: %w", err)
		}
	}

	newBalance := wallet.Balance
	if opType == model.OperationDeposit {
		newBalance += amountCents
	} else {
		if wallet.Balance < amountCents {
			return &model.OperationResult{
				WalletID:        walletID,
				BalanceCents:    wallet.Balance,
				OperationStatus: "INSUFFICIENT_BALANCE",
			}, model.ErrInsufficientFunds
		}
		newBalance -= amountCents
	}

	if err := txRepo.UpdateWalletBalance(ctx, walletID, newBalance); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	if err := txRepo.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &model.OperationResult{
		WalletID:        walletID,
		BalanceCents:    newBalance,
		OperationStatus: "SUCCESS",
	}, nil
}

func (u *WalletUseCase) GetBalance(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error) {
	return u.repo.GetWallet(ctx, walletID)
}
