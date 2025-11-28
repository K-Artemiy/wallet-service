package model

import (
	"errors"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrInsufficientFunds    = errors.New("insufficient balanse")
	ErrInvalidAmount        = errors.New("invalid amount")
)

type Wallet struct {
	WalletID uuid.UUID
	Balance  int64 //в копейках
}

type OperationResult struct {
	WalletID        uuid.UUID
	BalanceCents    int64
	OperationStatus string
}
