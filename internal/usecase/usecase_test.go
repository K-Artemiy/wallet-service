package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"wallet-service/internal/model"
	"wallet-service/internal/repository"
)

// Эмулируем репозиторий в памяти.
type mockRepo struct {
	wallets map[uuid.UUID]*model.Wallet
	inTx    bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		wallets: make(map[uuid.UUID]*model.Wallet),
		inTx:    false,
	}
}

func (m *mockRepo) BeginTransaction(ctx context.Context) (repository.WalletRepository, error) {
	m.inTx = true
	return m, nil
}

func (m *mockRepo) Commit(ctx context.Context) error {
	m.inTx = false
	return nil
}

func (m *mockRepo) Rollback(ctx context.Context) error {
	m.inTx = false
	return nil
}

func (m *mockRepo) GetWalletForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, errors.New("no rows in result set")
	}
	return &model.Wallet{WalletID: w.WalletID, Balance: w.Balance}, nil
}

func (m *mockRepo) CreateWallet(ctx context.Context, w *model.Wallet) error {
	m.wallets[w.WalletID] = &model.Wallet{WalletID: w.WalletID, Balance: w.Balance}
	return nil
}

func (m *mockRepo) UpdateWalletBalance(ctx context.Context, id uuid.UUID, newBalance int64) error {
	w, ok := m.wallets[id]
	if !ok {
		return errors.New("wallet not found")
	}
	w.Balance = newBalance
	return nil
}

func (m *mockRepo) GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	return &model.Wallet{WalletID: w.WalletID, Balance: w.Balance}, nil
}

func TestProcessOperation_Deposit(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	res, err := uc.ProcessOperation(ctx, id, model.OperationDeposit, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("result is nil")
	}
	if res.BalanceCents != 10000 {
		t.Fatalf("expected balance 10000, got %d", res.BalanceCents)
	}
	if res.OperationStatus != "SUCCESS" {
		t.Fatalf("expected status SUCCESS, got %s", res.OperationStatus)
	}
}

func TestProcessOperation_WithdrawSuccess(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	_, err := uc.ProcessOperation(ctx, id, model.OperationDeposit, 10000)
	if err != nil {
		t.Fatalf("unexpected error on deposit: %v", err)
	}

	res, err := uc.ProcessOperation(ctx, id, model.OperationWithdraw, 3000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("result is nil")
	}
	if res.BalanceCents != 7000 {
		t.Fatalf("expected balance 7000, got %d", res.BalanceCents)
	}
	if res.OperationStatus != "SUCCESS" {
		t.Fatalf("expected status SUCCESS, got %s", res.OperationStatus)
	}
}

func TestProcessOperation_WithdrawInsufficientBalance(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	_, err := uc.ProcessOperation(ctx, id, model.OperationDeposit, 2000)
	if err != nil {
		t.Fatalf("unexpected error on deposit: %v", err)
	}

	res, err := uc.ProcessOperation(ctx, id, model.OperationWithdraw, 5000)
	if !errors.Is(err, model.ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
	if res == nil {
		t.Fatalf("result is nil")
	}
	if res.OperationStatus != "INSUFFICIENT_BALANCE" {
		t.Fatalf("expected status INSUFFICIENT_BALANCE, got %s", res.OperationStatus)
	}
	if res.BalanceCents != 2000 {
		t.Fatalf("expected balance 2000 (unchanged), got %d", res.BalanceCents)
	}
}

func TestProcessOperation_InvalidOperationType(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	_, err := uc.ProcessOperation(ctx, id, model.OperationType("UNKNOWN"), 1000)
	if !errors.Is(err, model.ErrInvalidOperationType) {
		t.Fatalf("expected ErrInvalidOperationType, got %v", err)
	}
}

func TestProcessOperation_ZeroAmount(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	_, err := uc.ProcessOperation(ctx, id, model.OperationDeposit, 0)
	if !errors.Is(err, model.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestProcessOperation_NegativeAmount(t *testing.T) {
	repo := newMockRepo()
	uc := NewWalletUseCase(repo)
	ctx := context.Background()
	id := uuid.New()

	_, err := uc.ProcessOperation(ctx, id, model.OperationDeposit, -100)
	if !errors.Is(err, model.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}
