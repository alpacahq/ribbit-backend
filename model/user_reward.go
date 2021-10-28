package model

func init() {
	Register(&UserReward{})
}

type UserReward struct {
	Base
	ID                   int     `json:"id"`
	UserID               int     `json:"user_id"`
	JournalID            string  `json:"journal_id"`
	ReferredBy           int     `json:"referred_by"`
	RewardValue          float32 `json:"reward_value"`
	RewardType           string  `json:"reward_type"`
	RewardTransferStatus bool    `json:"reward_transfer_status"`
	ErrorResponse        string  `json:"error_response"`
}
