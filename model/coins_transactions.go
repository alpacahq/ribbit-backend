package model

func init() {
	Register(&CoinStatement{})
}

// Verification stores randomly generated tokens that can be redeemed
type CoinStatement struct {
	Base
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Coins  int    `json:"coins"`
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Status bool   `json:"status"`
}
