package repository

import (
	"context"

	"wallet-service/internal/model"

	"github.com/google/uuid"
)

type WalletRepository interface {
	BeginTransaction(ctx context.Context) (WalletRepository, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetWalletForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	CreateWallet(ctx context.Context, w *model.Wallet) error
	UpdateWalletBalance(ctx context.Context, id uuid.UUID, newBalance int64) error
	GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
}
