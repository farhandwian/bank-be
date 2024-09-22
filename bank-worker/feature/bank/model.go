package bank

type TransferEvent struct {
	Transfer              string `json:"transaction_id"`
	Amount                int    `json:"amount"`
	PhoneNumberOriginUser string `json:"phone_number_origin_user"`
	TargetUser            string `json:"target_user"`
	Remarks               string `json:"remarks"`
	CreatedAt             string `json:"created_at"`
}
