package model

func init() {
	Register(&Reward{})
}

type Reward struct {
	Base
	ID                   int     `json:"id"`
	PerAccountLimit      int     `json:"per_account_limit"`
	ReferralKycReward    float64 `json:"referral_kyc_reward"`
	ReferralSignupReward float64 `json:"referral_signup_reward"`
	ReferreKycReward     float64 `json:"referre_Kyc_reward"`
}
