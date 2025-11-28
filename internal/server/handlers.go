package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"wallet-service/internal/model"
	"wallet-service/internal/usecase"

	"github.com/google/uuid"
)

type Handlers struct {
	uc *usecase.WalletUseCase
}

func NewHandlers(uc *usecase.WalletUseCase) *Handlers {
	return &Handlers{uc: uc}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func parseAmountToCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, model.ErrInvalidAmount
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, model.ErrInvalidAmount
	}

	cents := math.Round(f * 100)
	return int64(cents), nil
}

func formatCentsToAmount(cents int64) string {
	rubles := float64(cents) / 100.0
	return strconv.FormatFloat(rubles, 'f', 2, 64)
}

func (h *Handlers) CreateOperationHandler(w http.ResponseWriter, r *http.Request) {
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		http.Error(w, "invalid walletId", http.StatusBadRequest)
		return
	}

	opType := model.OperationType(strings.ToUpper(req.OperationType))

	amountCents, err := parseAmountToCents(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount format", http.StatusBadRequest)
		return
	}

	res, err := h.uc.ProcessOperation(r.Context(), walletID, opType, amountCents)
	if err != nil {
		switch err {
		case model.ErrInsufficientFunds:
			writeJSON(w, http.StatusPaymentRequired, operationResponse{
				WalletID:        res.WalletID.String(),
				Balance:         formatCentsToAmount(res.BalanceCents),
				OperationStatus: res.OperationStatus,
				Error:           "not_enough_money",
			})
			return
		case model.ErrInvalidOperationType:
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		case model.ErrInvalidAmount:
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusOK, operationResponse{
		WalletID:        res.WalletID.String(),
		Balance:         formatCentsToAmount(res.BalanceCents),
		OperationStatus: res.OperationStatus,
	})
}

func (h *Handlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	walletIDStr := r.PathValue("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		http.Error(w, "invalid walletId", http.StatusBadRequest)
		return
	}

	wallet, err := h.uc.GetBalance(r.Context(), walletID)
	if err != nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, balanceResponse{
		WalletID: wallet.WalletID.String(),
		Balance:  formatCentsToAmount(wallet.Balance),
	})
}
