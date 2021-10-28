package model

func init() {
	Register(&BankAccount{})
}

// Verification stores randomly generated tokens that can be redeemed
type BankAccount struct {
	Base
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
	AccountID   string `json:"account_id"`
	BankName    string `json:"bank_name"`
	AccountName string `json:"account_name"`
	Status      bool   `json:"status"`
}
