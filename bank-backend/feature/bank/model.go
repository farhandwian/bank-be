package bank

type TopUpRequest struct {
	Amount int `json:"amount"validate:"required,min=1,numeric"`
}

type TopUpResponse struct {
	TopUpId       string `json:"top_up_id"`
	BalanceBefore int    `json:"balance_before"`
	BalanceAfter  int    `json:"balance_after"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks,omitempty"`
	CreatedAt     string `json:"created_at"`
}

type PaymentRequest struct {
	Amount  int    `json:"amount"validate:"required,min=1,numeric"`
	Remarks string `json:"remarks" validate:"required,max=50"`
}

type PaymentResponse struct {
	PaymentID     string `json:"payment_id"`
	BalanceBefore int    `json:"balance_before"`
	BalanceAfter  int    `json:"balance_after"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks,omitempty"`
	CreatedAt     string `json:"created_at"`
}

type TransferRequest struct {
	Amount     int    `json:"amount"validate:"required,min=1,numeric"`
	TargetUser string `json:"target_user"validate:"required"`
	Remarks    string `json:"remarks" validate:"required,max=50"`
}

type TransferResponse struct {
	PaymentID      string `json:"payment_id"`
	BalanceBefore  int    `json:"balance_before"`
	BalanceAfter   int    `json:"balance_after"`
	TargetTransfer string `json:"target_transfer"`
	Amount         int    `json:"amount"`
	Remarks        string `json:"remarks,omitempty"`
	CreatedAt      string `json:"created_at"`
}
