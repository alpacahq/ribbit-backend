package model

func init() {
	Register(&Verification{})
}

// Verification stores randomly generated tokens that can be redeemed
type Verification struct {
	Base
	ID     int    `json:"id"`
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}
