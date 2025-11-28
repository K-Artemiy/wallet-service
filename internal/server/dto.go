package server

type operationRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        string `json:"amount"`
}

type operationResponse struct {
	WalletID        string `json:"walletId"`
	Balance         string `json:"balance"`
	OperationStatus string `json:"operationStatus"`
	Error           string `json:"error,omitempty"`
}

type balanceResponse struct {
	WalletID string `json:"walletId"`
	Balance  string `json:"balance"`
}
