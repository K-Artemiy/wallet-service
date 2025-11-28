package pgrepo

import (
	"context"
	"errors"

	"wallet-service/internal/model"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletRepositoryPG struct {
	conn *pgx.Conn
	tx   pgx.Tx
}

func NewWalletRepository(conn *pgx.Conn) *WalletRepositoryPG {
	return &WalletRepositoryPG{conn: conn}
}

func (r *WalletRepositoryPG) checkTx() bool {
	return r.tx != nil
}

func (r *WalletRepositoryPG) BeginTransaction(ctx context.Context) (repository.WalletRepository, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	tr := &WalletRepositoryPG{tx: tx}
	return tr, nil
}

func (r *WalletRepositoryPG) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r *WalletRepositoryPG) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}

func (r *WalletRepositoryPG) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if r.checkTx() {
		return r.tx.QueryRow(ctx, sql, args...)
	}
	return r.conn.QueryRow(ctx, sql, args...)
}

func (r *WalletRepositoryPG) exec(ctx context.Context, sql string, args ...any) error {
	if r.checkTx() {
		_, err := r.tx.Exec(ctx, sql, args...)
		return err
	}
	_, err := r.conn.Exec(ctx, sql, args...)
	return err
}

func (r *WalletRepositoryPG) GetWalletForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	row := r.queryRow(ctx, `
		SELECT wallet_id, balance
		FROM wallets
		WHERE wallet_id = $1
		FOR UPDATE
	`, id)

	var w model.Wallet
	err := row.Scan(&w.WalletID, &w.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &w, nil
}

func (r *WalletRepositoryPG) CreateWallet(ctx context.Context, w *model.Wallet) error {
	return r.exec(ctx, `
		INSERT INTO wallets (wallet_id, balance)
		VALUES ($1, $2)
	`, w.WalletID, w.Balance)
}

func (r *WalletRepositoryPG) UpdateWalletBalance(ctx context.Context, id uuid.UUID, newBalance int64) error {
	return r.exec(ctx, `
		UPDATE wallets
		SET balance = $1
		WHERE wallet_id = $2
	`, newBalance, id)
}

func (r *WalletRepositoryPG) GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	row := r.queryRow(ctx, `
		SELECT wallet_id, balance
		FROM wallets
		WHERE wallet_id = $1
	`, id)

	var w model.Wallet
	if err := row.Scan(&w.WalletID, &w.Balance); err != nil {
		return nil, err
	}
	return &w, nil
}
